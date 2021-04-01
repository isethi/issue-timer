package main

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-github/github"
	"github.com/google/go-github/v33/github"
)

// integration test
func TestIssuesEvents(t *testing.T) {

	//set up

	githubToken := os.Getenv("INPUT_GITHUB_TOKEN")
	githubRepository := os.Getenv("GITHUB_REPOSITORY")

	ownerName, repoName := getOwnerAndRepoName(githubRepository)
	ctx := context.Background()
	client := getGitHubClient(ctx, githubToken)

	projectName := "1"
	labelName := "review"
	timeInColumnStr := "1"

	// create an issue

	input := &github.IssueRequest{
		Title:    github.String("t"),
		Body:     github.String("b"),
		Assignee: github.String("a"),
		Labels:   &[]string{"l1", "l2"},
	}

	iss, _, err := client.Issues.Create(ctx, ownerName, repoName, input)

	if err != nil {
		t.Errorf("Error creating issue: %v", err)
	}

	// get Project column details

	proj, _, err := client.Projects.GetProject(ctx, 1)
	if err != nil {
		t.Errorf("Error getting project: %v", err)
	}

	// add it to project board

	input = &github.ProjectCardOptions{
		ContentID:   iss.GetID(),
		ContentType: "Issue",
	}
	_, _, err = client.Projects.CreateProjectCard(ctx, 1, input)
	if err != nil {
		t.Errorf("Error creating project card: %v", err)
	}

	// run fn to make sure label is not added

	ParseIssues(projectName, githubToken, githubRepository, labelName, timeInColumnStr)

	iss, _, err = client.Issues.Get(ctx, ownerName, repoName, iss.GetNumber())
	if err != nil {
		t.Errorf("Error creating issue: %v", err)
	}
	// todo: check if label in issue labels

	// sleep for a minute

	// confirm label is added

	// move issue to another column

	// confirm label is removed

	//move back to old column

	// confirm label is not added back

	// delete issue

	// delete board
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
