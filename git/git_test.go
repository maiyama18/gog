package git

import (
	"github.com/stretchr/testify/assert"
	"io"
	"path"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		confFile       io.Reader
		expectedConfig *Config
		expectedErrMsg string
	}{
		{
			name:           "no setting",
			confFile:       strings.NewReader(""),
			expectedConfig: &DefaultConfig,
		},
		{
			name: "custom setting",
			confFile: strings.NewReader(`
[core]
        repositoryformatversion = 1
        filemode = false
        bare = true
        logallrefupdates = true
        ignorecase = true
        precomposeunicode = true
[remote "origin"]
        url = git@github.com:muiscript/gog.git
        fetch = +refs/heads/*:refs/remotes/origin/*
`),
			expectedConfig: &Config{
				RepositoryFormatVersion: 1,
				FileMode:                false,
				Bare:                    true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := NewConfig(test.confFile)
			assert.Equal(t, test.expectedConfig, config)
			if test.expectedErrMsg != "" {
				assert.EqualError(t, err, test.expectedErrMsg)
			}
		})
	}
}

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
