package github

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func NewCommentsService(c *Client) *CommentsService {
	return &CommentsService{c: c}
}

type CommentsService struct {
	c *Client
}

type Comment struct {
	ID        int64     `json:"id"`
	NodeID    string    `json:"node_id"`
	URL       string    `json:"url"`
	IssueURL  string    `json:"issue_url"`
	HtmlURL   string    `json:"html_url"`
	Body      string    `json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetCommentsParams struct {
	Org  string
	Repo string
}

func (s *CommentsService) GetComments(ctx context.Context, params GetCommentsParams) ([]*Comment, error) {
	allComments := make([]*Comment, 0)
	path := fmt.Sprintf("/repos/%s/%s/issues/comments", params.Org, params.Repo)

	for {
		comments, next, err := s.getCommentsFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)
		if next == "" {
			break
		}
		path = next
	}

	return allComments, nil
}

type GetIssueCommentsParams struct {
	Org         string
	Repo        string
	IssueNumber int64
}

func (s *CommentsService) GetIssueComments(ctx context.Context, params GetIssueCommentsParams) ([]*Comment, error) {
	allComments := make([]*Comment, 0)
	path := fmt.Sprintf("/repos/%s/%s/issues/%d/comments", params.Org, params.Repo, params.IssueNumber)

	for {
		comments, next, err := s.getCommentsFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)
		if next == "" {
			break
		}
		path = next
	}

	return allComments, nil
}

func (s *CommentsService) getCommentsFromPath(ctx context.Context, path string) ([]*Comment, string, error) {
	req, err := s.c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.c.do(req)
	if err != nil {
		return nil, "", err
	}

	var comments []*Comment
	if err := s.c.decode(resp, &comments); err != nil {
		return nil, "", err
	}

	links := s.c.parseLinkHeader(resp.Header.Get("Link"))

	return comments, links["next"], nil
}
