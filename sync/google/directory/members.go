package directory

import "google.golang.org/api/admin/directory/v1"

type Member struct {
	Id     string `json:"id,omitempty"`
	Email  string `json:"email,omitempty"`
	Etag   string `json:"etag,omitempty"`
	Role   string `json:"role,omitempty"`
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

type MemberType struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

func (c *Service) retrieveMembers(groupId string) (map[string]*Member, error) {
	completeMembers := map[string]*Member{}
	nextPageToken := ""

	var members *admin.Members
	var err error

	for {
		members, nextPageToken, err = c.memberCall(groupId, nextPageToken)
		if err != nil {
			return nil, err
		}
		for _, k := range members.Members {
			completeMembers[k.Id] = toMember(k)
		}

		if nextPageToken == "" {
			break
		}
	}

	return completeMembers, nil
}

func toMember(member *admin.Member) *Member {
	return &Member{
		Id:     member.Id,
		Email:  member.Email,
		Etag:   member.Etag,
		Role:   member.Role,
		Status: member.Status,
		Type:   member.Type,
	}
}

func (c *Service) memberCall(groupId string, pageToken string) (*admin.Members, string, error) {
	listCall := c.directoryService.Members.List(groupId).MaxResults(10000)
	if pageToken != "" {
		listCall = listCall.PageToken(pageToken)
	}

	members, err := listCall.Do()
	if err != nil {
		return nil, "", err
	}

	return members, members.NextPageToken, nil
}
