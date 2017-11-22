package directory

import "google.golang.org/api/admin/directory/v1"

type Group struct {
	Id          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Email       string             `json:"email,omitempty"`
	ETag        string             `json:"etag,omitempty"`
	Aliases     []string           `json:"aliases,omitempty"`
	Members     map[string]*Member `json:"members,omitempty"`
}

func (c *Service) retrieveGroups() (map[string]*Group, error) {
	completeGroups := map[string]*Group{}
	nextPageToken := ""

	var groups *admin.Groups
	var err error

	for {
		groups, nextPageToken, err = c.groupCall(nextPageToken)
		if err != nil {
			return nil, err
		}
		for _, k := range groups.Groups {
			completeGroups[k.Id] = toGroup(k)
		}

		if nextPageToken == "" {
			break
		}
	}

	return completeGroups, nil
}

func toGroup(group *admin.Group) *Group {
	return &Group{
		Id:          group.Id,
		Name:        group.Name,
		Description: group.Description,
		Email:       group.Email,
		ETag:        group.Etag,
		Aliases:     group.Aliases,
	}
}

func (c *Service) groupCall(pageToken string) (*admin.Groups, string, error) {
	listCall := c.directoryService.Groups.List().Customer(c.customerId).MaxResults(10000)
	if pageToken != "" {
		listCall = listCall.PageToken(pageToken)
	}
	if c.domain != "" {
		listCall = listCall.Domain(c.domain)
	}

	groups, err := listCall.Do()
	if err != nil {
		return nil, "", err
	}

	return groups, groups.NextPageToken, nil
}
