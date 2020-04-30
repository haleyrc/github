package github_test

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/haleyrc/github"
)

func TestGetIssues(t *testing.T) {
	c := github.NewClient()
	// c.Debug = true
	c.Token = os.Getenv("GITHUB_TOKEN")

	issues, err := c.GetIssues("frazercomputing", "frazerdms")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(issues))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	enc.Encode(issues)
}

func isEpic(i *github.Issue) bool {
	for _, label := range i.Labels {
		if strings.ToLower(label.Name) == "epic" {
			return true
		}
	}
	return false
}
