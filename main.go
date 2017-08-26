package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func main() {
	const issueTitle = "[Issue]"
	const prTitle = "[PR]"

	var title string

	ctx := context.Background()
	client := github.NewClient(nil)
	repos := []string{"minitest-retry", "travel_base"}

	for _, repo_name := range repos {
		fmt.Printf("**** %s ****\n", repo_name)

		issues, _, err := client.Issues.ListByRepo(ctx, "y-yagi", repo_name, nil)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		if len(issues) > 0 {
			for _, issue := range issues {
				if issue.PullRequestLinks == nil {
					title = issueTitle
				} else {
					title = prTitle
				}

				fmt.Printf("%s %s: %s\n", title, issue.GetTitle(), issue.GetHTMLURL())
			}
			fmt.Printf("\n")
		}
	}
}
