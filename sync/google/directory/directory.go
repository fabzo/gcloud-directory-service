package directory

import (
	"net/http"

	"google.golang.org/api/admin/directory/v1"
)

const (
	GroupType = "GROUP"
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

func ToMemberIdGroupIdsMapping(groups map[string]*Group) map[string][]string {
	memberIdToGroupIds := make(map[string]map[string]struct{})
	for _, group := range groups {
		groupId := group.Id
		for _, member := range group.Members {
			addMemberGroup(memberIdToGroupIds, member.Id, groupId)
		}
	}
	result := map[string][]string{}
	for memberId, groupIdsMap := range memberIdToGroupIds {
		groups := make([]string, 0)
		for groupId := range groupIdsMap {
			groups = append(groups, groupId)
		}
		result[memberId] = groups
	}
	return result
}

func addMemberGroup(memberIdToGroupIds map[string]map[string]struct{}, memberId string, groupId string) {
	if _, ok := memberIdToGroupIds[memberId]; !ok {
		memberIdToGroupIds[memberId] = make(map[string]struct{})
	}
	memberIdToGroupIds[memberId][groupId] = struct{}{}
}

func ToEmailMemberMapping(groups map[string]*Group) map[string]MemberType {
	emails := map[string]MemberType{}
	for id, group := range groups {
		emails[group.Email] = MemberType{
			Id:   id,
			Type: GroupType,
		}
		for _, alias := range group.Aliases {
			emails[alias] = MemberType{
				Id:   id,
				Type: GroupType,
			}
		}
		for _, member := range group.Members {
			emails[member.Email] = MemberType{
				Id:   id,
				Type: member.Type,
			}
		}
	}
	return emails
}
