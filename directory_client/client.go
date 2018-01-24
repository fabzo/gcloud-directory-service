package directory_client

import (
	"net/http"
	"net/url"

	"encoding/json"
	"fmt"
	"github.com/fabzo/gcloud-directory-service/sync/google/directory"
	"io/ioutil"
	"time"
)

type Client struct {
	url      string
	username string
	password string

	httpClient *http.Client

	groups             map[string]*directory.Group
	memberIdToGroupIds map[string][]string
	emailToMember      map[string]directory.MemberType
}

type cookieJar struct {
	jar map[string][]*http.Cookie
}

func (p *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *cookieJar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func New(url string, username string, password string) *Client {
	return &Client{
		url:      url,
		username: username,
		password: password,
		httpClient: &http.Client{
			Jar: &cookieJar{
				jar: make(map[string][]*http.Cookie),
			},
			Timeout: time.Second * 15,
		},
	}
}

func (c *Client) SyncDirectory() error {
	req, err := http.NewRequest("GET", c.url+"/api/directory", nil)
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("directory service request failed %d: %s", resp.StatusCode, string(data))
	}

	var groups map[string]*directory.Group
	err = json.Unmarshal(data, &groups)
	if err != nil {
		return err
	}

	c.groups = groups
	c.emailToMember = directory.ToEmailMemberMapping(groups)
	c.memberIdToGroupIds = directory.ToMemberIdGroupIdsMapping(groups)

	return nil
}

func (c *Client) Directory() map[string]*directory.Group {
	return c.groups
}

func (c *Client) MemberIdToGroupIdsMapping() map[string][]string {
	return c.memberIdToGroupIds
}

func (c *Client) EmailToMemberMapping() map[string]directory.MemberType {
	return c.emailToMember
}
