package github

import "fmt"

type Issue struct {
	ID               int       `json:"id"`
	NodeID           string    `json:"node_id"`
	URL              string    `json:"url"`
	RepositoryURL    string    `json:"repository_url"`
	LabelsURL        string    `json:"labels_url"`
	CommentsURL      string    `json:"comments_url"`
	EventsURL        string    `json:"events_url"`
	HtmlURL          string    `json:"html_url"`
	Number           int       `json:"number"`
	State            string    `json:"state"`
	Title            string    `json:"title"`
	Body             string    `json:"body"`
	User             User      `json:"user"`
	Labels           []Label   `json:"labels"`
	Assignee         User      `json:"assignee"`
	Assignees        []User    `json:"assignees"`
	Milestone        Milestone `json:"milestone"`
	Locked           bool      `json:"locked"`
	ActiveLockReason string    `json:"active_lock_reason"`
	Comments         int64     `json:"comments"`
	PullRequest      struct {
		URL      string `json:"url"`
		HtmlURL  string `json:"html_url"`
		DiffURL  string `json:"diff_url"`
		PatchURL string `json:"patch_url"`
	} `json:"pull_request"`
	ClosedAt  *string `json:"closed_at"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	ClosedBy  User    `json:"closed_by"`
}

func (c *Client) GetIssues(org, repo string) ([]*Issue, error) {
	path := fmt.Sprintf("%s/repos/%s/%s/issues", BaseURL, org, repo)

	var issues []*Issue
	if err := c.get(path, issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *Client) GetIssue(org, repo, number string) (*Issue, error) {
	path := fmt.Sprintf("%s/repos/%s/%s/issues/%s", BaseURL, org, repo, number)

	var issue Issue
	if err := c.get(path, &issue); err != nil {
		return nil, err
	}

	return &issue, nil
}
