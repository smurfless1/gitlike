package worktree

import (
	"github.com/sirupsen/logrus"
	"github.com/smurfless1/gitlike/repo"
	"github.com/smurfless1/pathlib"
	"github.com/stretchr/testify/assert"
	"testing"
)

var GitLocation string = pathlib.New("~/src/go/fancyrun/.git").ExpandUser().String()

func TestWorktree_Clone(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	var branch string

	cloneRoot := pathlib.New("/tmp/dt-cloneroot")
	if cloneRoot.Exists() {
		assert.Nil(t, cloneRoot.RmDir())
	}
	rootRepo := repo.New(GitLocation, cloneRoot, "main")
	treeRoot := pathlib.New("/tmp/dt-worktree-clone") // snrk, tree root
	if treeRoot.Exists() {
		assert.Nil(t, treeRoot.RmDir())
	}

	tree := New(GitLocation, cloneRoot, treeRoot, "main")
	assert.False(t, tree.Exists())
	assert.Nil(t, err)
	err = tree.Mkdir()
	err = tree.Clone()
	assert.Nil(t, err)
	branch, err = rootRepo.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "not-main", branch)
	branch, err = tree.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "main", branch)
}

func TestWorktree_CreateBranch(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	var branch string

	cloneRoot := pathlib.New("/tmp/dt-cloneroot")
	if cloneRoot.Exists() {
		assert.Nil(t, cloneRoot.RmDir())
	}
	rootRepo := repo.New(GitLocation, cloneRoot, "main")
	treeRoot := pathlib.New("/tmp/dt-worktree-create") // snrk, tree root
	if treeRoot.Exists() {
		assert.Nil(t, treeRoot.RmDir())
	}

	tree := New(GitLocation, cloneRoot, treeRoot, "fake-branch")
	assert.False(t, tree.Exists())
	assert.Nil(t, err)
	err = tree.Mkdir()
	err = tree.CreateBranch()
	assert.Nil(t, err)
	branch, err = rootRepo.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "not-main", branch)
	branch, err = tree.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "fake-branch", branch)
}
