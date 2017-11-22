package directory

import (
	"net/http"

	"google.golang.org/api/admin/directory/v1"
)

type Service struct {
	directoryService *admin.Service
	customerId       string
	domain           string
}

func New(client *http.Client, customerId string, domain string) (*Service, error) {
	service, err := admin.New(client)
	if err != nil {
		return nil, err
	}

	return &Service{
		directoryService: service,
		customerId:       customerId,
		domain:           domain,
	}, nil
}

func (c *Service) RetrieveDirectory() (map[string]*Group, error) {
	groups, err := c.retrieveGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		members, err := c.retrieveMembers(group.Id)
		if err != nil {
			return nil, err
		}
		group.Members = members
	}

	return groups, nil
}

func ToMemberGroupMapping(groups map[string]*Group) map[string][]string {
	members := map[string][]string{}

	for _, group := range groups {
		updateMemberMap(members, groups, group, group.Id)
	}

	return members
}

func updateMemberMap(members map[string][]string, groups map[string]*Group, group *Group, groupId string) {
	for memberId, _ := range group.Members {
		if _, ok := groups[memberId]; ok {
			updateMemberMap(members, groups, groups[memberId], groupId)
		} else {
			if groups, ok := members[memberId]; ok {
				members[memberId] = appendUnique(groups, groupId)
			} else {
				members[memberId] = []string{groupId}
			}
		}
	}
}

func appendUnique(groupIds []string, groupId string) []string {
	for _, id := range groupIds {
		if id == groupId {
			return groupIds
		}
	}
	return append(groupIds, groupId)
}

func ToEmailGroupMapping(groups map[string]*Group) map[string]string {
	emails := map[string]string{}
	for id, group := range groups {
		emails[group.Email] = id
		for _, alias := range group.Aliases {
			emails[alias] = id
		}
	}
	return emails
}
