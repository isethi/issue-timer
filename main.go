package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("In action")
	projectName := os.Getenv("INPUT_PROJECT_ID")
	githubToken := os.Getenv("INPUT_GITHUB_TOKEN")
	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	labelName := os.Getenv("INPUT_LABEL_NAME")
	timeInColumnStr := os.Getenv("INPUT_TIME_IN_MINUTES")

	ParseIssues(projectName, githubToken, githubRepository, labelName, timeInColumnStr)
	return
}
