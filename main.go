package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/github"
)

type config struct {
	User         string   `toml:"user"`
	Repositories []string `toml:"repositories"`
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

func run() int {
	const issueTitle = "[Issue]"
	const prTitle = "[PR]"

	var title string

	var cfg config
	err := cfg.load()
	if err != nil {
		return msg(err)
	}

	ctx := context.Background()
	client := github.NewClient(nil)

	for _, repo_name := range cfg.Repositories {
		fmt.Printf("**** %s ****\n", repo_name)

		issues, _, err := client.Issues.ListByRepo(ctx, cfg.User, repo_name, nil)
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

	return 0
}

func main() {
	os.Exit(run())
}
