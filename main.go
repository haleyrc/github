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

func (c *Client) Login(user, pass, id, secret string) error {
	url := fmt.Sprintf("%s/authorizations/clients/%s", BaseURL, id)

	requestBody := struct {
		Secret string   `json:"client_secret"`
		Scopes []string `json:"scopes"`
	}{secret, []string{"repo"}}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(requestBody); err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", url, &buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, pass)

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

	var response = struct {
		Token string `json:"token"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	c.Token = response.Token
	fmt.Println(c.Token)

	return nil
}

func (c *Client) GetIssues(org, repo string) ([]*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues", BaseURL, org, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Message)
	}

	var issues []*Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}

	return issues, nil
}

func (c *Client) GetIssue(org, repo, number string) (*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%s", BaseURL, org, repo, number)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Message)
	}

	var issue Issue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

func (c *Client) CreatePullRequestFromIssue(org, repo, head string, number int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", BaseURL, org, repo)

	requestBody := struct {
		Issue int    `json:"issue"`
		Head  string `json:"head"`
		Base  string `json:"base"`
	}{number, head, "master"}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(requestBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, err
		}
		msg := errResp.Message
		for _, e := range errResp.Errors {
			msg += "; " + e.Message
		}
		return nil, errors.New(errResp.Message)
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func dump(resp *http.Response) {
	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
}