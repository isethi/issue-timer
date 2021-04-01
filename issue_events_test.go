package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// integration test
func TestIssuesEvents(t *testing.T) {

	//set up

	githubToken := os.Getenv("GITHUB_TOKEN")
	githubRepository := os.Getenv("GITHUB_REPOSITORY")

	ownerName, repoName := getOwnerAndRepoName(githubRepository)
	ctx := context.Background()
	client := getGitHubClient(ctx, githubToken)

	projectName := "1"
	projectID := int64(1)
	labelName := "review"
	timeInColumnStr := "1"

	fmt.Printf("projectName is %s\n", projectName)
	fmt.Printf("github token is %s\n", githubToken)
	fmt.Printf("Github Repository is %s\n", githubRepository)
	fmt.Printf("Time in minutes is %s\n", timeInColumnStr)
	fmt.Printf("Label Name is %s\n", labelName)

	// create an issue

	input := &github.IssueRequest{
		Title:  github.String("t"),
		Body:   github.String("b"),
		Labels: &[]string{"l1", "l2"},
	}

	iss, _, err := client.Issues.Create(ctx, ownerName, repoName, input)

	if err != nil {
		t.Errorf("Error creating issue: %v", err)
	}

	// get Project column details

	projColumns, _, err := client.Projects.ListProjectColumns(ctx, projectID, &github.ListOptions{})
	if err != nil {
		t.Errorf("Error getting projec columnst: %v", err)
	}

	// add it to project board

	projectInput := &github.ProjectCardOptions{
		ContentID:   iss.GetID(),
		ContentType: "Issue",
	}
	projCard, _, err := client.Projects.CreateProjectCard(ctx, projectID, projectInput)
	if err != nil {
		t.Errorf("Error creating project card: %v", err)
	}
	_, err = client.Projects.MoveProjectCard(ctx, projCard.GetID(), &github.ProjectCardMoveOptions{ColumnID: projColumns[0].GetID()})
	if err != nil {
		t.Errorf("Error moving project card: %v", err)
	}
	// run fn to make sure label is not added

	ParseIssues(projectName, githubToken, githubRepository, labelName, timeInColumnStr)

	iss, _, err = client.Issues.Get(ctx, ownerName, repoName, iss.GetNumber())
	if err != nil {
		t.Errorf("Error creating issue: %v", err)
	}

	labelExists := checkLabelInIssueLabels(iss.Labels, labelName)
	if labelExists {
		t.Fatalf("Label should not exist but exists")
	}

	// sleep for a minute
	time.Sleep(time.Minute * 1)

	// confirm label is added
	labelExists = checkLabelInIssueLabels(iss.Labels, labelName)
	if !labelExists {
		t.Fatalf("Label should  exist but does not exist")
	}
	// move issue to another column

	_, err = client.Projects.MoveProjectCard(ctx, projCard.GetID(), &github.ProjectCardMoveOptions{ColumnID: projColumns[1].GetID()})
	if err != nil {
		t.Errorf("Error moving project card: %v", err)
	}
	// confirm label is removed
	ParseIssues(projectName, githubToken, githubRepository, labelName, timeInColumnStr)

	labelExists = checkLabelInIssueLabels(iss.Labels, labelName)
	if labelExists {
		t.Fatalf("Label should not exist but exists")
	}
	//move back to old column
	_, err = client.Projects.MoveProjectCard(ctx, projCard.GetID(), &github.ProjectCardMoveOptions{ColumnID: projColumns[0].GetID()})
	if err != nil {
		t.Errorf("Error moving project card: %v", err)
	}

	// confirm label is not added back

	ParseIssues(projectName, githubToken, githubRepository, labelName, timeInColumnStr)

	labelExists = checkLabelInIssueLabels(iss.Labels, labelName)
	if labelExists {
		t.Fatalf("Label should not exist but exists")
	}

	// delete issue - no client api to delete issue
	client.Issues.Edit(ctx, ownerName, repoName, iss.GetNumber(), &github.IssueRequest{State: github.String("closed")})
	// delete board - no client api to delete board

}

func checkLabelInIssueLabels(list []*github.Label, labelName string) bool {
	for _, b := range list {
		if b.GetName() == labelName {
			return true
		}
	}
	return false
}

func getGitHubClient2(ctx context.Context, githubToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}
