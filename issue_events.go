package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

// ParseIssues parses issues in a project and adds labels to issues excedding time in their current column
func ParseIssues(projectName, githubToken, githubRepository, labelName, timeInColumnStr string) {

	timeInColumn, err := strconv.ParseFloat(timeInColumnStr, 64)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("projectName is %s\n", projectName)
	fmt.Printf("github token is %s\n", githubToken)
	fmt.Printf("Github Repository is %s\n", githubRepository)
	fmt.Printf("Time in minutes is %f\n", timeInColumn)
	fmt.Printf("Label Name is %s\n", labelName)

	ownerName, repoName := getOwnerAndRepoName(githubRepository)
	ctx := context.Background()
	client := getGitHubClient(ctx, githubToken)

	// list all issues for the repository
	issues, _, err := client.Issues.ListByRepo(ctx, ownerName, repoName, &github.IssueListByRepoOptions{})
	if err != nil {
		fmt.Printf("Could not list issues in repo. Error: %v\n", err)
		return
	}
	fmt.Printf("issues are %v\n", issues)

	// for each issue list issue events
	for _, issue := range issues {
		events, _, err := client.Issues.ListIssueEvents(ctx, ownerName, repoName, *issue.Number, &github.ListOptions{})
		if err != nil {
			fmt.Printf("Could not list events in issue. Error: %v\n", err)
			return
		}
		fmt.Printf("Event: %v\n", events)

		for i := len(events) - 1; i >= 0; i-- {
			event := events[i]
			fmt.Printf("Event %d: %s, created: %s\n", i, *event.Event, event.CreatedAt)
			if event.ProjectCard != nil {
				fmt.Printf("project column: %s\n", *event.ProjectCard.ColumnName)
				diff := time.Now().Sub(*event.CreatedAt)
				fmt.Println(diff)
				if diff.Minutes() > timeInColumn {

					_, _, err := client.Issues.AddLabelsToIssue(ctx, ownerName, repoName, *issue.Number, []string{labelName})
					if err != nil {
						fmt.Printf("Could not add labels in issue. Error: %v\n", err)
						return
					}
				} else {
					_, err := client.Issues.RemoveLabelForIssue(ctx, ownerName, repoName, *issue.Number, labelName)
					if err != nil {
						fmt.Printf("Could not remove labels in issue. Error: %v\n", err)
						return
					}
				}
				break // if we found that an issue has been in a column then we dont parse older events
			}
		}
	}
	return
}

func getGitHubClient(ctx context.Context, githubToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func getOwnerAndRepoName(githubRepository string) (string, string) {
	githubRepositoryArray := strings.Split(githubRepository, "/")
	ownerName := githubRepositoryArray[0]
	repoName := githubRepositoryArray[1]
	return ownerName, repoName
}