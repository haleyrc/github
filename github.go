package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

// TODO (RCH): Parse errors
// {"message":"Validation Failed","errors":[{"resource":"PullRequest","code":"custom","message":"No commits between master and 20190222-test-issue-3"}],"documentation_url":"https://developer.github.com/v3/pulls/#create-a-pull-request"}

const BaseURL = "https://api.github.com"

type Issue struct {
	ID     int      `json:"id"`
	Number int      `json:"number"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []*Label `json:"labels"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type PullRequest struct {
	URL string `json:"html_url"`
}

type ErrorResponse struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	Errors           []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func NewClient() *Client {
	return &Client{
		hc: &http.Client{Timeout: 5 * time.Second},
	}
}

type Client struct {
	Token string
	hc    *http.Client
}

func (c *Client) MustLogin(user, pass, id, secret string) {
	if err := c.Login(user, pass, id, secret); err != nil {
		panic(err)
	}
}

func (c *Client) makeRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	url := BaseURL + path
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}

	return req, nil
}

func (c *Client) do(req *http.Request, data interface{}) error {
	resp, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return err
		}
		return errors.New(errResp.Message)
	}

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return err
	}

	return nil
}

func (c *Client) get(path string, result interface{}) error {
	req, err := c.makeRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	if err := c.do(req, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	req, err := c.makeRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}

	if err := c.do(req, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) Login(user, pass, id, secret string) error {
	path := fmt.Sprintf("/authorizations/clients/%s", id)

	requestBody := struct {
		Secret string   `json:"client_secret"`
		Scopes []string `json:"scopes"`
	}{
		Secret: secret,
		Scopes: []string{"repo"},
	}

	req, err := c.makeRequest(http.MethodPut, path, requestBody)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)

	var response struct {
		Token string `json:"token"`
	}
	if err := c.do(req, &response); err != nil {
		return err
	}

	c.Token = response.Token

	return nil
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

func (c *Client) CreatePullRequestFromIssue(org, repo, head string, number int) (*PullRequest, error) {
	path := fmt.Sprintf("%s/repos/%s/%s/pulls", BaseURL, org, repo)

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

func dump(resp *http.Response) {
	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
}
