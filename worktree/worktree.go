package worktree

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/smurfless1/fancyrun"
	git "github.com/smurfless1/gitlike"
	"github.com/smurfless1/gitlike/repo"
	"github.com/smurfless1/pathlib"
)

type Worktree struct {
	git.RepoLike
	Repo repo.Repo
	Base git.RepoBase
}

func New(remote string, clonedroot pathlib.Path, treeroot pathlib.Path, branch string) Worktree {
	out := Worktree{
		Repo: repo.New(remote, clonedroot, branch),
		Base: git.RepoBase{Branch: branch, Root: treeroot},
	}
	return out
}

func (w Worktree) Clone() error {
	// this requires more than just the base clone
	logrus.Debug("Worktree cloning")
	var err error
	if !w.Repo.Exists() {
		err = w.Repo.Clone()
		fancyrun.CheckInline(err)
		err = w.Repo.SetBranch("not-main")
		fancyrun.CheckInline(err)
	}
	err = w.Base.Mkdir()
	fancyrun.CheckInline(err)

	_, _, err = fancyrun.FancyRun(fmt.Sprintf("git worktree add %s %s", w.Base.Root.String(), w.Base.Branch), w.Repo.Base.Root, false)
	fancyrun.CheckInline(err)

	logrus.Debug("# Setting worktree branch")
	err = w.SetBranch(w.Base.Branch)
	return fancyrun.CheckFinal(err)
}

func (w Worktree) CreateBranch() error {
	// this requires more than just the base clone
	logrus.Debug("Worktree creating a new branch")
	var err error
	if !w.Repo.Exists() {
		err = w.Repo.Clone()
		fancyrun.CheckInline(err)
		err = w.Repo.SetBranch("not-main")
		fancyrun.CheckInline(err)
	}
	err = w.Base.Mkdir()
	fancyrun.CheckInline(err)

	_, _, err = fancyrun.FancyRun(fmt.Sprintf("git worktree add %s -b %s", w.Base.Root.String(), w.Base.Branch), w.Repo.Base.Root, false)
	fancyrun.CheckInline(err)

	logrus.Debug("# Creating new worktree branch")
	err = w.SetBranch(w.Base.Branch)
	return fancyrun.CheckFinal(err)
}

func (w Worktree) Reset() error {
	return w.Base.Reset()
}

func (w Worktree) Mkdir() error {
	err1 := w.Repo.Mkdir()
	err2 := w.Base.Mkdir()
	if !(err1 == nil && err2 == nil) {
		return errors.New("Failed to make a worktree directory")
	}
	parent, _ := w.Base.Root.Parent()
	if !parent.Exists() {
		err := parent.MkDir(os.FileMode(int(0o755)), true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w Worktree) Exists() bool {
	return w.Repo.Exists() && w.Base.Root.Exists() && w.Base.Root.JoinPath(".git").Exists()
}

func (w Worktree) Fetch() error {
	return w.Base.Fetch()
}

func (w Worktree) Pull() error {
	return w.Base.Pull()
}

func (w Worktree) SetBranch(value string) error {
	return w.Base.SetBranch(value)
}

func (w *Worktree) InitSubmodules() error {
	return w.Base.InitSubmodules()
}

func (w *Worktree) InitSubmodulesLazy(subdirPaths []pathlib.Path) error {
	return w.Base.InitSubmodulesLazy(subdirPaths)
}

func (w Worktree) CurrentGitHash() (string, error) {
	return w.Base.CurrentGitHash()
}

func (w Worktree) GitHashHistory() (string, error) {
	return w.Base.GitHashHistory()
}

func (w Worktree) ReadBranchFromGit() (string, error) {
	return w.Base.ReadBranchFromGit()
}

func (w Worktree) ReadRemoteFromGit(remote string) (string, error) {
	return w.Base.ReadRemoteFromGit(remote)
}
