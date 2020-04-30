package github

import "fmt"

type PullRequest struct {
	URL string `json:"html_url"`
}

func (c *Client) CreatePullRequestFromIssue(org, repo, head string, number int) (*PullRequest, error) {
	path := fmt.Sprintf("/repos/%s/%s/pulls", org, repo)

	requestBody := struct {
		Issue int    `json:"issue"`
		Head  string `json:"head"`
		Base  string `json:"base"`
	}{
		Issue: number,
		Head:  head,
		Base:  "master",
	}
	var pr PullRequest
	if err := c.post(path, requestBody, &pr); err != nil {
		return nil, err
	}

	return &pr, nil
}
