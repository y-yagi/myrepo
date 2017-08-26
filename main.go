package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

func main() {
	ctx := context.Background()

	client := github.NewClient(nil)
	repos := []string{"minitest-retry", "travel_base"}

	for _, repo_name := range repos {
		fmt.Printf("**** %s ****\n", repo_name)
		prs, _, err := client.PullRequests.List(ctx, "y-yagi", repo_name, nil)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		if len(prs) > 0 {
			fmt.Printf("PRs:\n")
			for _, pr := range prs {
				fmt.Printf("%s: %s\n", pr.GetTitle(), pr.GetHTMLURL())
			}
			fmt.Printf("\n")
		}

		issues, _, err := client.Issues.ListByRepo(ctx, "y-yagi", repo_name, nil)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			continue
		}

		if len(issues) > 0 {
			fmt.Printf("Issues:\n")
			for _, issue := range issues {
				if issue.PullRequestLinks == nil {
					fmt.Printf("%s: %s\n", issue.GetTitle(), issue.GetHTMLURL())
				}
			}
			fmt.Printf("\n")
		}
	}
}
