## GCloud Directory Service
[![Travis CI Build](https://travis-ci.org/fabzo/gcloud-directory-service.svg?branch=master "Travis CI Build")](https://travis-ci.org/fabzo/gcloud-directory-service)

Provides a REST accessible cache of google apps groups and members.

### Requirements

- A service account with domain wide delegate activated (service account)
- The email address of a user that has the required access rights (subject)

### Setup

Follow the guide for the G Suite Admin Directory API access [here](https://developers.google.com/admin-sdk/directory/v1/guides/delegation#create_the_service_account_and_its_credentials)

The required auth scopes are:

    https://www.googleapis.com/auth/admin.directory.group.readonly
    https://www.googleapis.com/auth/admin.directory.group.member.readonly

Additional to the service account with these permissions a actual user account in G Suites is required. The account needs to have the same access rights as above. Usually a domain admin account can be used for this purpose, as the service account cannot gain more permissions as given through the security settings. This ensures that the required permissions are available on the user side.

These are the minimum requirements to run the directory service.

Docker Example:

	docker run --rm -it \
		-v $(pwd):/account fabzo/gcloud-directory-service server \
		--service-account /account/service-account.json \
		--subject admin@your.org


### Building

	go build

### Help output

    Run the directory server

    Usage:
      gcloud-directory-service server [flags]

    Flags:
      -b, --basic-auth string         Basic auth login in the form of <username>:<password>. Random login is generated if not set.
      -c, --customer-id string        The gsuite customer id. Defaults to my_customer. (default "my_customer")
      -d, --domain string             The gsuite domain for which to retrieve the groups. Defaults to ''
      -h, --help                      help for server
      -p, --port int                  Port for the API (default: 8080) (default 8080)
      -a, --service-account string    Location of the service account json file
      -l, --storage-location string   Storage location for the directory for faster restores (optional)
      -s, --subject string            The gsuite user to impersonate
      -i, --sync-interval int         Sync interval in minutes. Defaults to 30. (default 30)


### Using the Go client library

There is a simple implementation of a client library in directory_client that does nothing more than to retrieve the entire directory.
It provides all information of the /api/directory, /api/groups and /api/members endpoints locally.

    func main() {
        client := directory_client.New("https://directory-service.url", "username", "password")
        // Call SyncDirectory to update the local directory copy
        err := client.SyncDirectory()
        if err != nil {
            fmt.Errorf("Failed: %v\n", err)
            os.Exit(1)
        }

        json, _ := json.Marshal(client.Directory())

        fmt.Printf("Directory: %s\n", json)
    }

### API endpoints:

    /
    /api
        Link list with endpoints

    /api/status
        Information about the sync status and directory content
        {
        	"last_sync": ...,
			"last_sync_duration": ...,
			"next_sync": ...,
			"known_groups": 0,
			"known_users": 0,
			"sync_in_progress": false
        }

    /api/directory
        The entire directory with group to member mappings
        {
			"cryptic group id 1": {
				"id": "cryptic group id 1",
				"name": "group name",
				"email": "group email",
				"etag": "etag",
				"members": {
					"cryptic user id 1": {
						"id": "cryptic user id 1",
						"email": "users email address",
						"etag": "etag",
						"role": "MEMBER",
						"status": "ACTIVE",
						"type": "USER"
				},
				...
			},
			...
		}

    /api/groups
        Mapping of group email addresses to group IDs
        {
			"somegroup1@your.org": "cryptic user id 1",
			"groupalias@your.org": "cryptic user id 1",
			"somegroup3@your.org": "cryptic user id 2",
			"somegroup4@your.org": "cryptic user id 3",
			...
        }

    /api/members
        Mapping of member IDs to group IDs they are part of
        {
			"cryptic user id 1": ["cryptic group id 1", "cryptic group id 2"],
			"cryptic user id 2": ["cryptic group id 1", "cryptic group id 3"],
			"cryptic user id 3": ["cryptic group id 2", "cryptic group id 4"],
			...
        }

    /health
        Always returns 200 OK