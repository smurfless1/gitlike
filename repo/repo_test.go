package repo

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/smurfless1/pathlib"
	"github.com/stretchr/testify/assert"
)

const GitLocation string = "/Users/davidb/src/go/gobyexample/.git"

func TestRepo_Clone(t *testing.T) {
	t.Skip()
	root := pathlib.New("/tmp/dt")
	repo := New(GitLocation, root, "master")
	err := repo.Clone()
	assert.Nil(t, err)
	assert.True(t, root.Exists())
	err = repo.Reset()
	assert.Nil(t, err)
	assert.True(t, root.Exists())
	var branch string
	branch, err = repo.ReadBranchFromGit()
	assert.Nil(t, err)
	assert.Equal(t, "master", branch)
}

func TestRepo_read_remote(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	root := pathlib.New("~/src/go/gobyexample").ExpandUser()
	repo := New(GitLocation, root, "master")
	remote, err := repo.ReadRemoteFromGit("origin")
	assert.Nil(t, err)
	assert.Equal(t, "git@github.com:smurfless1/gitlike.git", remote)
	remote, err = repo.ReadRemoteFromGit("missingremote")
	assert.Equal(t, "None", remote)
	assert.NotNil(t, err)
	assert.Equal(t, "unable to locate a remote with that name", err.Error())
}
