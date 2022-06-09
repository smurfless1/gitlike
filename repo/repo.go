package repo

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/smurfless1/fancyrun"
	"github.com/smurfless1/gitlike"
	"github.com/smurfless1/pathlib"
	"os"
)

type Repo struct {
	git.RepoLike
	Remote string
	Base   git.RepoBase
}

func (r *Repo) Clone() error {
	logrus.Debug("Repo cloning")
	err := r.Base.Root.MkDir(os.FileMode(int(0755)), true)
	if err != nil {
		return err
	}
	workdir, err := r.Base.Root.Parent()
	if err != nil {
		logrus.Error(err)
		return err
	}

	if r.Base.Root.JoinPath(".git").Exists() {
		return nil
	}
	_, _, err = fancyrun.FancyRun(
		fmt.Sprintf("git clone %s %s", r.Remote, r.Base.Root.String()),
		workdir,
		true,
	)
	return err
}

func (r *Repo) Reset() error {
	return r.Base.Reset()
}

func (r *Repo) Mkdir() error {
	return r.Base.Mkdir()
}

func (r *Repo) Exists() bool {
	return r.Base.Root.Exists() && r.Base.Root.JoinPath(".git").Exists()
}

func (r *Repo) Fetch() error {
	return r.Base.Fetch()
}

func (r *Repo) Pull() error {
	return r.Base.Pull()
}

func (r *Repo) ReadBranchFromGit() (string, error) {
	return r.Base.ReadBranchFromGit()
}

func (r *Repo) SetBranch(value string) error {
	return r.Base.SetBranch(value)
}

func (r *Repo) InitSubmodules() error {
	return r.Base.InitSubmodules()
}

func (r *Repo) InitSubmodulesLazy(subdirPaths []pathlib.Path) error {
	return r.Base.InitSubmodulesLazy(subdirPaths)
}

func (r *Repo) CurrentGitHash() (string, error) {
	return r.Base.CurrentGitHash()
}

func (r *Repo) GitHashHistory() (string, error) {
	return r.Base.GitHashHistory()
}

func New(remote string, Root pathlib.Path, branch string) Repo {
	out := Repo{
		Base: git.RepoBase{
			Root:   Root,
			Branch: branch,
		},
		Remote: remote,
	}
	return out
}
