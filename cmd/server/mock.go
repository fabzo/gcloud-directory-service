package server

import (
	"net/http"
	"os"
	"strings"

	"strconv"

	"github.com/fabzo/gcloud-directory-service/sync"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	Mock.PersistentFlags().StringVarP(&basicAuth, "basic-auth", "b", "", "Basic auth login in the form of <username>:<password>.")
	Mock.PersistentFlags().StringVarP(&storageLocation, "storage-location", "l", "", "Storage location where the directory.json is located")
	Mock.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Port for the API")
}

var Mock = &cobra.Command{
	Use:   "mock",
	Short: "Run the mock server",
	Run: func(cmd *cobra.Command, args []string) {

		if basicAuth != "" && !strings.Contains(basicAuth, ":") {
			logrus.Errorf("Missing colon in basic auth argument. Format is <username>:<password>.")
			os.Exit(1)
		}
		if basicAuth == "" {
			logrus.Errorf("No basic auth login provided.")
			os.Exit(1)
		}

		if storageLocation == "" {
			logrus.Errorf("No storage location provided.")
			os.Exit(1)
		}

		mockSync, err := sync.Mock(storageLocation)
		if err != nil {
			logrus.Errorf("Could not initiate mock client: %v", err)
			os.Exit(1)
		}
		logrus.Infof("Starting mock server")
		logrus.Infof("server port          : %v", port)
		logrus.Infof("basic auth           : %v", basicAuth)
		logrus.Infof("storage location     : %v", storageLocation)

		mockSync.RunSyncLoop()

		r := mux.NewRouter()
		r.HandleFunc("/", auth(rootHandler()))
		r.HandleFunc("/api", auth(rootHandler()))
		r.HandleFunc("/api/status", auth(statusHandler(mockSync)))
		r.HandleFunc("/api/directory", auth(directoryHandler(mockSync)))
		r.HandleFunc("/api/groups", auth(groupsHandler(mockSync)))
		r.HandleFunc("/api/members", auth(membersHandler(mockSync)))
		r.HandleFunc("/health", healthHandler())

		http.ListenAndServe(":"+strconv.Itoa(port), r)
	},
}
