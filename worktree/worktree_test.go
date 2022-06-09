package worktree

import (
	"github.com/sirupsen/logrus"
	"github.com/smurfless1/gitlike/repo"
	"github.com/smurfless1/pathlib"
	"github.com/stretchr/testify/assert"
	"testing"
)

var GitLocation string = pathlib.New("~/work/devicetests/.git").ExpandUser().String()

func TestRepo_Clone(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	var err error
	var branch string

	cloneRoot := pathlib.New("/tmp/dt-cloneroot")
	if cloneRoot.Exists() {
		assert.Nil(t, cloneRoot.RmDir())
	}
	rootRepo := repo.New(GitLocation, cloneRoot, "master")
	treeRoot := pathlib.New("/tmp/dt-worktree") // snrk, tree root
	if treeRoot.Exists() {
		assert.Nil(t, treeRoot.RmDir())
	}

	tree := New(GitLocation, cloneRoot, treeRoot, "master")
	assert.False(t, tree.Exists())
	assert.Nil(t, err)
	err = tree.Mkdir()
	err = tree.Clone()
	assert.Nil(t, err)
	branch, err = rootRepo.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "not_master", branch)
	branch, err = tree.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch)
	// todo the devicetests object would store these paths
	err = tree.InitSubmodulesLazy([]pathlib.Path{tree.Base.Root.JoinPath("common_sdk", ".git")})
	assert.Nil(t, err)
}
