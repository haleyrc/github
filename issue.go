package github

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Issue struct {
	ID               int64      `json:"id"`
	NodeID           string     `json:"node_id"`
	URL              string     `json:"url"`
	RepositoryURL    string     `json:"repository_url"`
	LabelsURL        string     `json:"labels_url"`
	CommentsURL      string     `json:"comments_url"`
	EventsURL        string     `json:"events_url"`
	HtmlURL          string     `json:"html_url"`
	Number           int64      `json:"number"`
	State            string     `json:"state"`
	Title            string     `json:"title"`
	Body             string     `json:"body"`
	User             *User      `json:"user"`
	Labels           []Label    `json:"labels"`
	Assignee         *User      `json:"assignee"`
	Assignees        []User     `json:"assignees"`
	Milestone        *Milestone `json:"milestone"`
	Locked           bool       `json:"locked"`
	ActiveLockReason string     `json:"active_lock_reason"`
	Comments         int64      `json:"comments"`
	PullRequest      *struct {
		URL      string `json:"url"`
		HtmlURL  string `json:"html_url"`
		DiffURL  string `json:"diff_url"`
		PatchURL string `json:"patch_url"`
	} `json:"pull_request"`
	ClosedAt  *string   `json:"closed_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClosedBy  *User     `json:"closed_by"`
}

func NewIssuesService(c *Client) *IssuesService {
	return &IssuesService{c: c}
}

type IssuesService struct {
	c *Client
}

type GetIssueParams struct {
	Org    string
	Repo   string
	Number string
}

func (s *IssuesService) GetIssue(ctx context.Context, params GetIssueParams) (*Issue, error) {
	path := fmt.Sprintf("/repos/%s/%s/issues/%s", params.Org, params.Repo, params.Number)

	var issue Issue
	if err := s.c.get(ctx, path, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

type GetIssuesParams struct {
	Org  string
	Repo string
}

func (s *IssuesService) GetIssues(ctx context.Context, params GetIssuesParams) ([]*Issue, error) {
	allIssues := make([]*Issue, 0)
	path := fmt.Sprintf("/repos/%s/%s/issues?status=open", params.Org, params.Repo)

	for {
		issues, next, err := s.getIssuesFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, issues...)
		if next == "" {
			break
		}
		path = next
	}

	return allIssues, nil
}

func (s *IssuesService) getIssuesFromPath(ctx context.Context, path string) ([]*Issue, string, error) {
	req, err := s.c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.c.do(req)
	if err != nil {
		return nil, "", err
	}

	var issues []*Issue
	if err := s.c.decode(resp, &issues); err != nil {
		return nil, "", err
	}

	links := s.c.parseLinkHeader(resp.Header.Get("Link"))

	return issues, links["next"], nil
}
