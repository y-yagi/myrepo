package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"github.com/y-yagi/configure"
)

type config struct {
	User         string   `toml:"user"`
	Repositories []string `toml:"repositories"`
	AccessToken  string   `toml:"access_token"`
}

func msg(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 0
}

func fetchFromGitHub(ctx context.Context, client *github.Client, user string, repo string, ch chan<- string) {
	var result string
	result += fmt.Sprintf("**** %s ****\n", repo)

	issues, _, err := client.Issues.ListByRepo(ctx, user, repo, nil)
	if err != nil {
		result += fmt.Sprintf("err: %v\n", err)
		ch <- result
		return
	}

	if len(issues) > 0 {
		for _, issue := range issues {
			if issue.PullRequestLinks == nil {
				result += "[Issue] "
			} else {
				result += "[PR] "
			}

			result += fmt.Sprintf("%s: %s\n", issue.GetTitle(), issue.GetHTMLURL())
		}
		result += fmt.Sprintf("\n")
	}

	ch <- result
}

func run() int {
	var editConfig bool
	var cfg config
	var client *github.Client

	flags := flag.NewFlagSet("myrepo", flag.ExitOnError)
	flags.BoolVar(&editConfig, "c", false, "Edit config.")
	flags.Parse(os.Args[1:])

	if editConfig {
		editor := os.Getenv("EDITOR")
		if len(editor) == 0 {
			editor = "vim"
		}

		if err := configure.Edit("myrepo", editor); err != nil {
			return msg(err)
		}

		return 0
	}

	err := configure.Load("myrepo", &cfg)
	if err != nil {
		return msg(err)
	}

	ch := make(chan string)

	ctx := context.Background()
	if len(cfg.AccessToken) > 0 {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.AccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	for _, repoName := range cfg.Repositories {
		go fetchFromGitHub(ctx, client, cfg.User, repoName, ch)
	}

	for range cfg.Repositories {
		fmt.Fprintf(os.Stdout, <-ch)
	}

	return 0
}

func main() {
	os.Exit(run())
}
