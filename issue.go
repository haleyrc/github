package github

import (
	"fmt"
	"net/http"
)

type Issue struct {
	ID               int        `json:"id"`
	NodeID           string     `json:"node_id"`
	URL              string     `json:"url"`
	RepositoryURL    string     `json:"repository_url"`
	LabelsURL        string     `json:"labels_url"`
	CommentsURL      string     `json:"comments_url"`
	EventsURL        string     `json:"events_url"`
	HtmlURL          string     `json:"html_url"`
	Number           int        `json:"number"`
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
	ClosedAt  *string `json:"closed_at"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ClosedBy  *User   `json:"closed_by"`
}

// TODO (RCH): Take query params
func (c *Client) GetIssues(org, repo string) ([]*Issue, error) {
	allIssues := make([]*Issue, 0)
	path := fmt.Sprintf("/repos/%s/%s/issues?status=open", org, repo)
	for {
		fmt.Println("Fetching", path)
		issues, next, err := c.getIssuesFromPath(path)
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

func (c *Client) getIssuesFromPath(path string) ([]*Issue, string, error) {
	req, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, "", err
	}

	var issues []*Issue
	if err := c.decode(resp, &issues); err != nil {
		return nil, "", err
	}

	links := c.parseLinkHeader(resp.Header.Get("Link"))

	return issues, links["next"], nil
}

// func (c *Client) GetIssues(org, repo string) ([]*Issue, error) {
// 	path := fmt.Sprintf("/repos/%s/%s/issues", org, repo)

// 	var issues []*Issue
// 	if err := c.get(path, &issues); err != nil {
// 		return nil, err
// 	}

// 	return issues, nil
// }

func (c *Client) GetIssue(org, repo, number string) (*Issue, error) {
	path := fmt.Sprintf("/repos/%s/%s/issues/%s", org, repo, number)

	var issue Issue
	if err := c.get(path, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}
