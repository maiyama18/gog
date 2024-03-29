package git

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"strings"
	"testing"
)

var nonExistentDirectory = "./testdata/nonexistent01"
var notGitRepository = "./testdata/not_repository01"
var gitRepository = "./testdata/repository01"

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
			expectedErrMsg: "not a git repository",
		},
		{
			name:           "failure - not git repository without force",
			workTree:       notGitRepository,
			force:          false,
			expectedErrMsg: "not a git repository",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo, err := NewRepository(test.workTree, test.force)

			if test.expectedErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.True(t, strings.Contains(err.Error(), test.expectedErrMsg))
			}
			assert.Equal(t, test.expectedRepository, repo)
		})
	}
}

const emptyDirectory = "./testdata/empty01"
const notEmptyDirectory = "./testdata/not_empty01"
const regularFile = "./testdata/regular_file01"

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
		{
			name:           "failure - path already exists as file",
			repositoryPath: regularFile,
			expectedErrMsg: "already exist",
			teardown:       func() {},
		},
		{
			name:           "failure - dir is not empty",
			repositoryPath: notEmptyDirectory,
			expectedErrMsg: "not empty",
			teardown:       func() {},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.teardown()

			repo, err := CreateRepository(test.repositoryPath)
			if test.expectedErrMsg == "" {
				assert.Nil(t, err)

				assert.Equal(t, test.expectedRepository, repo)

				assert.True(t, isExistingDir(test.repositoryPath))
				assert.True(t, isExistingDir(path.Join(test.repositoryPath, ".git")))
				assert.True(t, isExistingFile(path.Join(test.repositoryPath, ".git", "config")))
			} else {
				assert.True(t, strings.Contains(err.Error(), test.expectedErrMsg), fmt.Sprintf("expected '%s' to contain '%s'", err.Error(), test.expectedErrMsg))
			}
		})
	}
}

const shaHello = "ce013625030ba8dba906f756967f9e9ca394464a"

func TestRepository_ReadObject(t *testing.T) {
	tests := []struct {
		name           string
		sha            string
		kind           string
		expectedObject Object
		expectedErrMsg string
	}{
		{
			name:           "success",
			sha:            shaHello,
			kind:           "blob",
			expectedObject: NewBlob([]byte("hello\n")),
		},
		{
			name:           "failure - invalid sha",
			sha:            "invalid sha",
			kind:           "blob",
			expectedErrMsg: "not a directory",
		},
		{
			name:           "failure - wrong kind",
			sha:            shaHello,
			kind:           "commit",
			expectedErrMsg: "type mismatch",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r, err := NewRepository(gitRepository, true)
			if err != nil {
				t.Fatal(err)
				t.FailNow()
			}
			obj, err := r.ReadObject(test.sha, test.kind)
			if test.expectedErrMsg == "" {
				assert.Nil(t, err)

				assert.Equal(t, test.expectedObject, obj)
			} else {
				assert.True(t, strings.Contains(err.Error(), test.expectedErrMsg), fmt.Sprintf("expected '%s' to contain '%s'", err.Error(), test.expectedErrMsg))
			}
		})
	}
}
