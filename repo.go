package github

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Repo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Private     bool      `json:"private"`
	Description string    `json:"description"`
	Fork        bool      `json:"fork"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PushedAt    time.Time `json:"pushed_at"`
	Archived    bool      `json:"archived"`
	Disabled    bool      `json:"disabled"`
	HtmlURL     string    `json:"html_url"`
}

func NewReposService(c *Client) *ReposService {
	return &ReposService{c: c}
}

type ReposService struct {
	c *Client
}

type GetReposParams struct {
	Org string
}

func (s *ReposService) GetRepos(ctx context.Context, params GetReposParams) ([]*Repo, error) {
	allRepos := []*Repo{}
	path := fmt.Sprintf("/orgs/%s/repos", params.Org)

	for {
		repos, next, err := s.getReposFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)
		if next == "" {
			break
		}
		path = next
	}

	return allRepos, nil
}

func (s *ReposService) getReposFromPath(ctx context.Context, path string) ([]*Repo, string, error) {
	req, err := s.c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.c.do(req)
	if err != nil {
		return nil, "", err
	}

	var repos []*Repo
	if err := s.c.decode(resp, &repos); err != nil {
		return nil, "", err
	}

	links := s.c.parseLinkHeader(resp.Header.Get("Link"))

	return repos, links["next"], nil
}
