package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"golang.org/x/oauth2"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
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

func configDir() string {
	var dir string

	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "myrepo")
		}
		dir = filepath.Join(dir, "myrepo")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "myrepo")
	}
	return dir
}

func (cfg *config) load() error {
	dir := configDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("cannot create directory: %v", err)
	}
	file := filepath.Join(dir, "config.toml")

	_, err := os.Stat(file)
	if err == nil {
		_, err := toml.DecodeFile(file, cfg)
		if err != nil {
			return err
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	return nil
}

func fetchFromGitHub(client *github.Client, ctx context.Context, user string, repo string, ch chan<- string) {
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
	var client *github.Client

	var cfg config
	err := cfg.load()
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

	for _, repo_name := range cfg.Repositories {
		go fetchFromGitHub(client, ctx, cfg.User, repo_name, ch)
	}

	for range cfg.Repositories {
		fmt.Printf(<-ch)
	}

	return 0
}

func main() {
	os.Exit(run())
}
