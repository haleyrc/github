package github

import (
	"context"
	"fmt"
	"net/http"
)

type User struct {
	Login             string `json:"login"`
	ID                int64  `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HtmlURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

func NewUsersService(c *Client) *UsersService {
	return &UsersService{c: c}
}

type UsersService struct {
	c *Client
}

type GetAssigneesParams struct {
	Org  string
	Repo string
}

func (s *UsersService) GetAssignees(ctx context.Context, params GetAssigneesParams) ([]*User, error) {
	allAssignees := []*User{}
	path := fmt.Sprintf("/repos/%s/%s/assignees", params.Org, params.Repo)

	for {
		assignees, next, err := s.getAssigneesFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		allAssignees = append(allAssignees, assignees...)
		if next == "" {
			break
		}
		path = next
	}

	return allAssignees, nil
}

func (s *UsersService) getAssigneesFromPath(ctx context.Context, path string) ([]*User, string, error) {
	req, err := s.c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.c.do(req)
	if err != nil {
		return nil, "", err
	}

	var assignees []*User
	if err := s.c.decode(resp, &assignees); err != nil {
		return nil, "", err
	}

	links := s.c.parseLinkHeader(resp.Header.Get("Link"))

	return assignees, links["next"], nil
}
