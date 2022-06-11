package git

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/smurfless1/fancyrun"
	"github.com/smurfless1/pathlib"
	"os"
	"strings"
)

type RepoLike interface {
	Clone() error
	Reset() error
	Mkdir() error
	Exists() bool
	Fetch() error
	Pull() error
	SetBranch(value string) error
	ReadBranchFromGit() (string, error)
	InitSubmodules() error
	InitSubmodulesLazy(subdirPaths []pathlib.Path) error
	CurrentGitHash() (string, error)
	GitHashHistory() (string, error)
	ReadRemoteFromGit(string) (string, error)
}

type RepoBase struct {
	Root   pathlib.Path
	Branch string
}

func (r *RepoBase) InitSubmodules() error {
	_, _, err := fancyrun.FancyRun("git submodule update --init --recursive", r.Root, true)
	return fancyrun.CheckFinal(err)
}

func (r *RepoBase) InitSubmodulesLazy(subdirPaths []pathlib.Path) error {
	if !AllExist(subdirPaths) {
		return r.InitSubmodules()
	}
	return nil
}
func (r *RepoBase) CurrentGitHash() (string, error) {
	_, out, err := fancyrun.FancyRun("git rev-parse HEAD", r.Root, true)
	readOut := strings.TrimRight(string(out), "\n")
	return readOut, fancyrun.CheckFinal(err)
}

func (r *RepoBase) GitHashHistory() (string, error) {
	_, out, err := fancyrun.FancyRun("git log --graph --pretty=%H", r.Root, true)
	readOut := strings.TrimRight(string(out), "\n")
	return readOut, fancyrun.CheckFinal(err)
}

func (r *RepoBase) ReadBranchFromGit() (string, error) {
	_, out, err := fancyrun.FancyRun("git branch --show-current", r.Root, true)
	if err != nil {
		return "", err
	}
	readOut := strings.TrimRight(string(out), "\n")
	r.Branch = readOut
	return r.Branch, nil
}

func (r *RepoBase) Mkdir() error {
	if !r.Root.Exists() {
		err := r.Root.MkDir(os.FileMode(int(0755)), true)
		if err != nil {
			return err
		}
	}
	return nil
}

func AllExist(vs []pathlib.Path) bool {
	for _, v := range vs {
		if !v.Exists() {
			return false
		}
	}
	return true
}

func (r *RepoBase) SetBranch(value string) error {
	_, _, err := fancyrun.FancyRun(fmt.Sprintf("git switch -c %s --discard-changes --recurse-submodules", value), r.Root, false)
	if err != nil {
		_, _, err = fancyrun.FancyRun(fmt.Sprintf("git switch %s --discard-changes --recurse-submodules", value), r.Root, true)
		if err != nil {
			logrus.Fatal("Unable to switch branches")
		}
	}
	r.Branch = value
	return nil
}

func (r *RepoBase) Pull() error {
	cmd, _, err := fancyrun.FancyRun("git pull", r.Root, false)
	// log_file_name="GitRepo.pull"
	if err != nil || cmd.ProcessState.ExitCode() != 0 {
		return r.Reset()
	}
	return nil
}

func (r *RepoBase) Fetch() error {
	_, _, err := fancyrun.FancyRunWithNamedLog("git fetch", r.Root, true, "git-fetch")
	return fancyrun.CheckFinal(err)
}

func (r *RepoBase) Reset() error {
	commands := []string{
		fmt.Sprintf("git reset --hard origin/%s", r.Branch),
		"git clean -fxd",
		fmt.Sprintf("git switch %s", r.Branch),
		"git clean -fxd",
		fmt.Sprintf("git switch %s", r.Branch),
		"git submodule foreach --recursive git clean -xfd",
		"git submodule foreach --recursive git reset --hard",
	}
	var err error
	for _, command := range commands {
		_, _, err = fancyrun.FancyRun(command, r.Root, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RepoBase) ReadRemoteFromGit(remote string) (string, error) {
	var err error
	var output string
	var pwd string
	pwd, err = os.Getwd()
	fancyrun.CheckInline(err)
	_, output, err = fancyrun.FancyRun("git remote -v", pathlib.New(pwd), false)
	fancyrun.CheckInline(err)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		if words[0] == remote {
			return words[1], nil
		}
	}
	return "None", errors.New("unable to locate a remote with that name")

}
