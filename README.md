# GCloud Directory Service

Provides a REST accessible cache of google apps groups and members.

Builds:

    https://travis-ci.org/fabzo/gcloud-directory-service

Docker:
    
    docker pull fabzo/gcloud-directory-service
    docker run --rm -it fabzo/gcloud-directory-service server --help


server --help output:

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


API endpoints:

    /
    /api
        Link list with endpoints

    /api/status
        Information about the sync status and directory content

    /api/directory
        The entire directory with group to member mappings

    /api/groups
        Mapping of group email addresses to group IDs

    /api/members
        Mapping of member IDs to group IDs they are part of

    /health
        Always returns 200 OK