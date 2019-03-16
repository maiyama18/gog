package git

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

var nonExistentDirectory = "./git/testdata/nonexistent01"
var notGitRepository = "./git/testdata/not_repository01"
var gitRepository = "./git/testdata/repository01"

func TestNewRepository(t *testing.T) {
	tests := []struct {
		name               string
		workTree           string
		force              bool
		expectedRepository *Repository
		expectedErrMsg     string
	}{
		{
			name:     "success - nonexistent directory",
			workTree: nonExistentDirectory,
			force:    true,
			expectedRepository: &Repository{
				WorkTree: nonExistentDirectory,
				GitDir:   path.Join(nonExistentDirectory, ".git"),
				Conf:     &defaultConfig,
			},
		},
		{
			name:     "success - not git repository",
			workTree: notGitRepository,
			force:    true,
			expectedRepository: &Repository{
				WorkTree: notGitRepository,
				GitDir:   path.Join(notGitRepository, ".git"),
				Conf:     &defaultConfig,
			},
		},
		{
			name:     "success - git repository",
			workTree: gitRepository,
			force:    true,
			expectedRepository: &Repository{
				WorkTree: gitRepository,
				GitDir:   path.Join(gitRepository, ".git"),
				Conf:     &defaultConfig,
			},
		},
		{
			name:           "failure - nonexistent directory without force",
			workTree:       nonExistentDirectory,
			force:          false,
			expectedErrMsg: "not a git repository: git/testdata/nonexistent01/.git",
		},
		{
			name:           "failure - not git repository without force",
			workTree:       notGitRepository,
			force:          false,
			expectedErrMsg: "not a git repository: git/testdata/not_repository01/.git",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo, err := newRepository(test.workTree, test.force)

			if test.expectedErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expectedErrMsg)
			}
			assert.Equal(t, test.expectedRepository, repo)
		})
	}
}

const emptyDirectory = "./git/testdata/empty01"

func TestCreateRepository(t *testing.T) {
	tests := []struct {
		name               string
		repositoryPath     string
		expectedRepository *Repository
		expectedErrMsg     string
		teardown           func()
	}{
		{
			name:           "success - empty dir",
			repositoryPath: emptyDirectory,
			expectedRepository: &Repository{
				WorkTree: emptyDirectory,
				GitDir:   path.Join(emptyDirectory, ".git"),
				Conf:     &defaultConfig,
			},
			teardown: func() {},
		},
		{
			name:           "success - nonexistent dir",
			repositoryPath: nonExistentDirectory,
			expectedRepository: &Repository{
				WorkTree: nonExistentDirectory,
				GitDir:   path.Join(nonExistentDirectory, ".git"),
				Conf:     &defaultConfig,
			},
			teardown: func() {
				_ = os.RemoveAll(nonExistentDirectory)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.teardown()

			repo, err := CreateRepository(test.repositoryPath)
			if test.expectedErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expectedErrMsg)
			}
			assert.Equal(t, test.expectedRepository, repo)

			assert.True(t, isExistingDir(test.repositoryPath))
			assert.True(t, isExistingDir(path.Join(test.repositoryPath, ".git")))
			assert.True(t, isExistingFile(path.Join(test.repositoryPath, ".git", "config")))
		})
	}
}
