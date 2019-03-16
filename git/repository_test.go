package git

import (
	"github.com/stretchr/testify/assert"
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
				Conf:     &DefaultConfig,
			},
		},
		{
			name:     "success - not git repository",
			workTree: notGitRepository,
			force:    true,
			expectedRepository: &Repository{
				WorkTree: notGitRepository,
				GitDir:   path.Join(notGitRepository, ".git"),
				Conf:     &DefaultConfig,
			},
		},
		{
			name:     "success - git repository",
			workTree: gitRepository,
			force:    true,
			expectedRepository: &Repository{
				WorkTree: gitRepository,
				GitDir:   path.Join(gitRepository, ".git"),
				Conf:     &DefaultConfig,
			},
		},
		{
			name:           "success - nonexistent directory withour force",
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
			repo, err := NewRepository(test.workTree, test.force)

			if test.expectedErrMsg == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.expectedErrMsg)
			}
			assert.Equal(t, test.expectedRepository, repo)
		})
	}
}
