package sync

import (
	"encoding/json"
	"errors"
	"github.com/fabzo/gcloud-directory-service/sync/google/directory"
	"io/ioutil"
	"os"
	"path/filepath"
)

type mockSync struct {
	groups             map[string]*directory.Group
	memberIdToGroupIds map[string][]string
	emailToMember      map[string]directory.MemberType
}

func Mock(storageLocation string) (DirSync, error) {
	groups, err := getGroupsFromDisk(storageLocation)
	if err != nil {
		return nil, err
	}
	return &mockSync{groups: groups,
		emailToMember:      directory.ToEmailMemberMapping(groups),
		memberIdToGroupIds: directory.ToMemberIdGroupIdsMapping(groups)}, nil
}

func getGroupsFromDisk(location string) (map[string]*directory.Group, error) {
	if location == "" {
		return nil, errors.New("storage location is empty")
	}

	fileInfo, err := os.Stat(location)
	if err != nil {
		return nil, err
	}

	file := location
	if fileInfo.IsDir() {
		file = filepath.Join(location, "directory.json")
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var groups map[string]*directory.Group
	err = json.Unmarshal(data, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (m *mockSync) RunSyncLoop() {
}

func (m *mockSync) Status() *Status {
	knownGroups := len(m.groups)
	userCounter := 0
	for _, group := range m.groups {
		userCounter += len(group.Members)
	}
	knownUsers := userCounter
	return &Status{KnownGroups: knownGroups, KnownUsers: knownUsers}
}

func (m *mockSync) Directory() map[string]*directory.Group {
	return m.groups
}

func (m *mockSync) MemberIdToGroupIdsMapping() map[string][]string {
	return m.memberIdToGroupIds
}

func (m *mockSync) EmailToMemberMapping() map[string]directory.MemberType {
	return m.emailToMember
}
