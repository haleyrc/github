package github_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/haleyrc/github"
)

func TestGetRepositories(t *testing.T) {
	c := github.NewClient()
	c.Token = os.Getenv("GITHUB_TOKEN")
	fmt.Println(c.Token)

	repos, err := c.ListRepositories("frazercomputing")
	if err != nil {
		t.Fatalf("failed to get repos: %v", err)
	}

	for _, repo := range repos {
		t.Logf("%s: %t", repo.Name, repo.Private)
	}
}
