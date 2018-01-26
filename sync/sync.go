package sync

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/fabzo/gcloud-directory-service/sync/google"
	"github.com/fabzo/gcloud-directory-service/sync/google/directory"
	"github.com/sirupsen/logrus"
)

type DirSync interface {
	RunSyncLoop()
	Status() *Status
	Directory() map[string]*directory.Group
	MemberIdToGroupIdsMapping() map[string][]string
	EmailToMemberMapping() map[string]directory.MemberType
}

type dirSync struct {
	serviceAccountFile string
	subject            string
	customerId         string
	domain             string
	syncInterval       int
	storageLocation    string

	syncRunningMutex sync.Mutex
	syncRunning      bool

	googleClient *google.Client

	groups             map[string]*directory.Group
	memberIdToGroupIds map[string][]string
	emailToMember      map[string]directory.MemberType

	status *Status
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

type Status struct {
	LastSync         time.Time `json:"last_sync"`
	LastSyncDuration Duration  `json:"last_sync_duration"`
	NextSync         time.Time `json:"next_sync"`
	KnownGroups      int       `json:"known_groups"`
	KnownUsers       int       `json:"known_users"`
	SyncInProgress   bool      `json:"sync_in_progress"`
}

func New(serviceAccountFile string, subject string, customerId string, domain string, syncInterval int, storageLocation string) (DirSync, error) {

	if serviceAccountFile == "" {
		return nil, fmt.Errorf("service account location cannot be empty")
	}
	if customerId == "" {
		return nil, fmt.Errorf("customer id cannot be empty")
	}
	if syncInterval < 5 {
		return nil, fmt.Errorf("sync interval cannot be lower than 5 minutes")
	}

	dirSync := &dirSync{
		serviceAccountFile: serviceAccountFile,
		subject:            subject,
		customerId:         customerId,
		domain:             domain,
		syncInterval:       syncInterval,
		storageLocation:    storageLocation,
		status:             &Status{},
		syncRunning:        false,
	}

	err := dirSync.restoreFromDisk(storageLocation)
	if err != nil {
		logrus.Warnf("Failed to restore directory from disk: %v", err)
	}

	return dirSync, nil
}

func (d *dirSync) RunSyncLoop() {
	d.syncRunningMutex.Lock()
	defer d.syncRunningMutex.Unlock()
	if !d.syncRunning {
		d.syncRunning = true
		go d.syncLoop()
	}
}

func (d *dirSync) syncLoop() {
	for true {

		if d.googleClient == nil {
			serviceAccount, err := ioutil.ReadFile(d.serviceAccountFile)
			if err != nil {
				logrus.Errorf("Could not read service account file. Skipping current sync attempt. Error: %v", err)
				goto skip
			}

			d.googleClient, err = google.New(serviceAccount, d.subject, d.customerId, d.domain)
			if err != nil {
				logrus.Errorf("Could not initiate google client. Skipping current sync attempt. Error: %v", err)
				goto skip
			}
		}

		d.executeSync()

	skip:
		time.Sleep(time.Duration(d.syncInterval) * time.Minute)
	}
}

func (d *dirSync) executeSync() {
	d.status.LastSync = time.Now()
	d.status.SyncInProgress = true

	groups, err := d.googleClient.Directory.RetrieveDirectory()
	if err != nil {
		logrus.Errorf("Failed to execute sync. Error: %v", err)
	} else {
		d.updateGroups(groups)

		err = d.persistToDisk(d.storageLocation)
		if err != nil {
			logrus.Warnf("Failed to persist directory to disk: %v", err)
		}
	}

	d.status.SyncInProgress = false
	d.status.LastSyncDuration = Duration{time.Since(d.status.LastSync)}
	d.status.NextSync = time.Now().Add(time.Duration(d.syncInterval) * time.Minute)
}

func (d *dirSync) updateStatusCounter(groups map[string]*directory.Group) {
	d.status.KnownGroups = len(d.groups)
	userCounter := 0
	for _, group := range d.groups {
		userCounter += len(group.Members)
	}
	d.status.KnownUsers = userCounter
}

func (d *dirSync) updateGroups(groups map[string]*directory.Group) {
	d.groups = groups

	d.emailToMember = directory.ToEmailMemberMapping(groups)
	d.memberIdToGroupIds = directory.ToMemberIdGroupIdsMapping(groups)
	d.updateStatusCounter(groups)
}

func (d *dirSync) Status() *Status {
	return d.status
}

func (d *dirSync) Directory() map[string]*directory.Group {
	return d.groups
}

func (d *dirSync) MemberIdToGroupIdsMapping() map[string][]string {
	return d.memberIdToGroupIds
}

func (d *dirSync) EmailToMemberMapping() map[string]directory.MemberType {
	return d.emailToMember
}

func (d *dirSync) persistToDisk(location string) error {
	if location == "" {
		return nil
	}

	data, err := json.Marshal(d.groups)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(location+"/directory.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (d *dirSync) restoreFromDisk(location string) error {
	if location == "" {
		return nil
	}

	data, err := ioutil.ReadFile(location + "/directory.json")
	if err != nil {
		return err
	}

	var groups map[string]*directory.Group
	err = json.Unmarshal(data, &groups)
	if err != nil {
		return err
	}

	d.updateGroups(groups)
	return nil
}
