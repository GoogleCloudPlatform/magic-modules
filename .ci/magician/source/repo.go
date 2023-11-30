package source

import (
	"fmt"
	"magician/exec"
	"path/filepath"
)

type Repo struct {
	Name        string // Name in GitHub (e.g. magic-modules)
	Title       string // Title for display (e.g. Magic Modules)
	Branch      string // Branch to clone, optional
	Path        string // local Path once cloned, including Name
	DiffCanFail bool   // whether to allow the command to continue if cloning or diffing the repo fails
}

type Controller struct {
	rnr      *exec.Runner
	username string
	token    string
	goPath   string
}

func NewController(goPath, username, token string, rnr *exec.Runner) *Controller {
	return &Controller{
		rnr:      rnr,
		username: username,
		token:    token,
		goPath:   goPath,
	}
}

func (gc Controller) SetPath(repo *Repo) {
	repo.Path = filepath.Join(gc.goPath, "src", "github.com", gc.username, repo.Name)
}

func (gc Controller) Clone(repo *Repo) error {
	var err error
	url := fmt.Sprintf("https://%s:%s@github.com/%s/%s", gc.username, gc.token, gc.username, repo.Name)
	if repo.Branch == "" {
		_, err = gc.rnr.Run("git", []string{"clone", url, repo.Path}, nil)
	} else {
		_, err = gc.rnr.Run("git", []string{"clone", "-b", repo.Branch, url, repo.Path}, nil)
	}
	return err
}

func (gc Controller) Fetch(repo *Repo, branch string) error {
	if err := gc.rnr.PushDir(repo.Path); err != nil {
		return err
	}
	if _, err := gc.rnr.Run("git", []string{"fetch", "origin", branch}, nil); err != nil {
		return fmt.Errorf("error fetching branch %s in repo %s: %v\n", branch, repo.Name, err)
	}
	return gc.rnr.PopDir()
}

func (gc Controller) Diff(repo *Repo, oldBranch, newBranch string) (string, error) {
	if err := gc.rnr.PushDir(repo.Path); err != nil {
		return "", err
	}
	diffs, err := gc.rnr.Run("git", []string{"diff", "origin/" + oldBranch, "origin/" + newBranch, "--shortstat"}, nil)
	if err != nil {
		return "", fmt.Errorf("error diffing %s and %s: %v", oldBranch, newBranch, err)
	}
	return diffs, gc.rnr.PopDir()
}
