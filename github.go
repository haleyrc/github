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

const DefaultURL = "https://api.github.com"

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		URL:        DefaultURL,
	}
}

type Client struct {
	Debug      bool
	HTTPClient *http.Client
	Token      string
	URL        string
}

func (c *Client) MustLogin(user, pass, id, secret string) {
	if err := c.Login(user, pass, id, secret); err != nil {
		panic(err)
	}
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

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := c.decode(resp, &response); err != nil {
		return err
	}

	c.Token = response.Token

	return nil
}

func (c *Client) makeRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	url := c.URL + path
	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if c.Token != "" {
		req.Header.Set("Authorization", "token "+c.Token)
	}

	if c.Debug {
		dumpRequest(req)
	}

	return req, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		dumpResponse(resp)
	}

	return resp, nil
}

func (c *Client) decode(resp *http.Response, data interface{}) error {
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

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if err := c.decode(resp, result); err != nil {
		return err
	}

	return nil
}

func (c *Client) post(path string, body interface{}, result interface{}) error {
	req, err := c.makeRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if err := c.decode(resp, result); err != nil {
		return err
	}

	return nil
}

func dumpResponse(resp *http.Response) {
	b, _ := httputil.DumpResponse(resp, true)
	fmt.Println(string(b))
}

func dumpRequest(req *http.Request) {
	b, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(b))
}
