// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/github"
	"github.com/hashicorp/go-changelog"
	"github.com/manifoldco/promptui"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"
)

//go:embed changelog-entry.tmpl
var changelogTmplDefault string

type Note struct {
	// service or area of codebase the pull request changes
	Subcategory string
	// release note type (Bug...)
	Type string
	// release note text
	Description string
	// pull request number
	PR int
	// URL of the pull request
	URL string
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var subcategory, changeType, description, changelogTmpl, dir, url string
	var pr int
	var Url bool
	flag.BoolVar(&Url, "add-url", false, "add GitHub issue URL (omitted by default due to formatting in changelog-build)")
	flag.IntVar(&pr, "pr", -1, "pull request number")
	flag.StringVar(&subcategory, "subcategory", "", "the service or area of the codebase the pull request changes (optional)")
	flag.StringVar(&changeType, "type", "", "the type of change")
	flag.StringVar(&description, "description", "", "the changelog entry content")
	flag.StringVar(&changelogTmpl, "changelog-template", "", "the path of the file holding the template to use for the changelog entries")
	flag.StringVar(&dir, "dir", "", "the relative path from the current directory of where the changelog entry file should be written")
	flag.Parse()

	if pr == -1 {
		pr, url, err = getPrNumberFromGithub(pwd)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Must specify pull request number or run in a git repo with a GitHub remote origin:", err)
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
			os.Exit(1)
		}
	}
	fmt.Fprintln(os.Stderr, "Found matching pull request:", url)

	if changeType == "" {
		prompt := promptui.Select{
			Label: "Select a change type",
			Items: changelog.TypeValues,
		}

		_, changeType, err = prompt.Run()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Must specify the change type")
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
			os.Exit(1)
		}
	} else {
		if !changelog.TypeValid(changeType) {
			fmt.Fprintln(os.Stderr, "Must specify a valid type")
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
			os.Exit(1)
		}
	}

	if subcategory == "" {
		prompt := promptui.Prompt{Label: "Subcategory (optional)"}
		subcategory, err = prompt.Run()
	}

	if description == "" {
		prompt := promptui.Prompt{Label: "Description"}
		description, err = prompt.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Must specify the change description")
			fmt.Fprintln(os.Stderr, "")
			flag.Usage()
			os.Exit(1)
		}
	}

	var tmpl *template.Template
	if changelogTmpl != "" {
		file, err := os.ReadFile(changelogTmpl)
		if err != nil {
			os.Exit(1)
		}
		tmpl, err = template.New("").Parse(string(file))
		if err != nil {
			os.Exit(1)
		}
	} else {
		tmpl, err = template.New("").Parse(changelogTmplDefault)
		if err != nil {
			os.Exit(1)
		}
	}

	if !Url {
		url = ""
	}
	n := Note{Type: changeType, Description: description, Subcategory: subcategory, PR: pr, URL: url}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, n)
	fmt.Printf("\n%s\n", buf.String())
	if err != nil {
		os.Exit(1)
	}
	filename := fmt.Sprintf("%d.txt", pr)
	filepath := path.Join(pwd, dir, filename)
	err = os.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "Created changelog entry at", filepath)
}

func OpenGit(path string) (*git.Repository, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		if path == "/" {
			return r, err
		} else {
			return OpenGit(path[:strings.LastIndex(path, "/")])
		}
	}
	return r, err
}

func getPrNumberFromGithub(path string) (int, string, error) {
	r, err := OpenGit(path)
	if err != nil {
		return -1, "", err
	}

	ref, err := r.Head()
	if err != nil {
		return -1, "", err
	}

	localBranch, err := r.Branch(ref.Name().Short())
	if err != nil {
		return -1, "", err
	}

	remote, err := r.Remote("origin")
	if err != nil {
		return -1, "", err
	}

	if len(remote.Config().URLs) <= 0 {
		return -1, "", errors.New("not able to parse repo and org")
	}
	remoteUrl := remote.Config().URLs[0]

	re := regexp.MustCompile(`.*github\.com:(.*)/(.*)\.git`)
	m := re.FindStringSubmatch(remoteUrl)
	if len(m) < 3 {
		return -1, "", errors.New("not able to parse repo and org")
	}

	cli := github.NewClient(nil)

	ctx := context.Background()

	githubOrg := m[1]
	githubRepo := m[2]

	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 200},
		Sort:        "updated",
		Direction:   "desc",
	}

	list, _, err := cli.PullRequests.List(ctx, githubOrg, githubRepo, opt)
	if err != nil {
		return -1, "", err
	}

	for _, pr := range list {
		head := pr.GetHead()
		if head == nil {
			continue
		}

		branch := head.GetRef()
		if branch == "" {
			continue
		}

		repo := head.GetRepo()
		if repo == nil {
			continue
		}

		// Allow finding PRs from forks - localBranch.Remote will return the
		// remote name for branches of origin, but the remote URL for forks
		var gitRemote string
		remote, err := r.Remote(localBranch.Remote)
		if err != nil {
			gitRemote = localBranch.Remote
		} else {
			gitRemote = remote.Config().URLs[0]
		}

		if (gitRemote == *repo.SSHURL || gitRemote == *repo.CloneURL) &&
			localBranch.Name == branch {
			n := pr.GetNumber()

			if n != 0 {
				return n, pr.GetHTMLURL(), nil
			}
		}
	}

	return -1, "", errors.New("no match found")
}
