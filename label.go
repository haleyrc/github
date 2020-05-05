package github

import (
	"context"
	"fmt"
	"net/http"
)

type Label struct {
	ID          int64  `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

func NewLabelsService(c *Client) *LabelsService {
	return &LabelsService{c: c}
}

type LabelsService struct {
	c *Client
}

type GetLabelsParams struct {
	Org  string
	Repo string
}

func (s *LabelsService) GetLabels(ctx context.Context, params GetLabelsParams) ([]*Label, error) {
	allLabels := []*Label{}
	path := fmt.Sprintf("/repos/%s/%s/labels", params.Org, params.Repo)

	for {
		labels, next, err := s.getLabelsFromPath(ctx, path)
		if err != nil {
			return nil, err
		}

		allLabels = append(allLabels, labels...)
		if next == "" {
			break
		}
		path = next
	}

	return allLabels, nil
}

func (s *LabelsService) getLabelsFromPath(ctx context.Context, path string) ([]*Label, string, error) {
	req, err := s.c.makeRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := s.c.do(req)
	if err != nil {
		return nil, "", err
	}

	var labels []*Label
	if err := s.c.decode(resp, &labels); err != nil {
		return nil, "", err
	}

	links := s.c.parseLinkHeader(resp.Header.Get("link"))

	return labels, links["next"], nil
}
