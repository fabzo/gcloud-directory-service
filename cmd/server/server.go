package server

import (
	"net/http"
	"fmt"
	"os"
	"encoding/json"
	"strings"

	"github.com/spf13/cobra"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/fabzo/gcloud-directory-service/sync"
	"github.com/fabzo/gcloud-directory-service/utils"
	"strconv"
)

var serviceAccount string
var subject string
var customerId string
var domain string
var syncInterval int
var storageLocation string
var port int

var basicAuth string

func init() {
	Command.PersistentFlags().StringVarP(&serviceAccount, "service-account", "a", "", "Location of the service account json file")
	Command.PersistentFlags().StringVarP(&subject, "subject", "s", "", "The gsuite user to impersonate")
	Command.PersistentFlags().StringVarP(&customerId, "customer-id", "c", "my_customer", "The gsuite customer id")
	Command.PersistentFlags().StringVarP(&domain, "domain", "d", "", "The gsuite domain for which to retrieve the groups (default '')")
	Command.PersistentFlags().IntVarP(&syncInterval, "sync-interval", "i", 30, "Sync interval in minutes")
	Command.PersistentFlags().StringVarP(&basicAuth, "basic-auth", "b", "", "Basic auth login in the form of <username>:<password>. Random login is generated if not set")
	Command.PersistentFlags().StringVarP(&storageLocation, "storage-location", "l", "", "Storage location for the directory for faster restores (optional)")
	Command.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port for the API")
}

var Command = &cobra.Command{
	Use:   "server",
	Short: "Run the directory server",
	Run: func(cmd *cobra.Command, args []string) {

		if basicAuth != "" && !strings.Contains(basicAuth, ":") {
			logrus.Errorf("Missing colon in basic auth argument. Format is <username>:<password>.")
			os.Exit(1)
		}
		if basicAuth == "" {
			basicAuth = "admin:" + utils.RandString(25)
			logrus.Warnf("No basic auth login provided. Randomly generated basic auth is " + basicAuth)
		}

		dirSync, err := sync.New(serviceAccount, subject, customerId, domain, syncInterval, storageLocation)
		if err != nil {
			logrus.Errorf("Could not initiate google sync client: %v", err)
			os.Exit(1)
		}

		dirSync.RunSyncLoop()

		r := mux.NewRouter()
		r.HandleFunc("/", auth(rootHandler()))
		r.HandleFunc("/status", auth(statusHandler(dirSync)))
		r.HandleFunc("/directory", auth(directoryHandler(dirSync)))
		r.HandleFunc("/groups", auth(groupsHandler(dirSync)))
		r.HandleFunc("/members", auth(membersHandler(dirSync)))
		http.ListenAndServe(":" + strconv.Itoa(port), r)

	},
}

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		user, pass, _ := r.BasicAuth()
		if !check(user, pass) {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}

func check(username string, password string) bool {
	return username + ":" + password == basicAuth;
}

func rootHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
<a href="/">/</a></br>
<a href="/status">/status</a></br>
<a href="/directory">/directory</a></br>
<a href="/groups">/groups</a></br>
<a href="/members">/members</a>
		`))
	}
}

func statusHandler(dirSync *sync.DirSync) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(dirSync.Status())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Failed to marshal status json: %v\n", err)))
			return
		}
	}
}

func directoryHandler(dirSync *sync.DirSync) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		groups := dirSync.Directory()
		if groups == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		} else {
			err := json.NewEncoder(w).Encode(groups)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Failed to marshal directory json: %v\n", err)))
				return
			}
		}
	}
}

func groupsHandler(dirSync *sync.DirSync) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		groups := dirSync.MailToGroupMapping()
		if groups == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		} else {
			err := json.NewEncoder(w).Encode(groups)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Failed to marshal groups json: %v\n", err)))
				return
			}
		}
	}
}

func membersHandler(dirSync *sync.DirSync) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		members := dirSync.MemberToGroupMapping()
		if members == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		} else {
			err := json.NewEncoder(w).Encode(members)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Failed to marshal members json: %v\n", err)))
				return
			}
		}
	}
}