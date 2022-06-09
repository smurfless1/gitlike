package repo

import (
	"github.com/smurfless1/pathlib"
	"github.com/stretchr/testify/assert"
	"testing"
)

// const GitLocation string = "https://stash.prod.netflix.net:7006/scm/nrdp/devicetests.git"
const GitLocation string = "/Users/davidb/work/devicetests/.git"

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
