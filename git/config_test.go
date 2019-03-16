package git

import (
	"github.com/stretchr/testify/assert"
	"io"
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
			expectedConfig: &defaultConfig,
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
			config, err := newConfig(test.confFile)
			assert.Equal(t, test.expectedConfig, config)
			if test.expectedErrMsg != "" {
				assert.EqualError(t, err, test.expectedErrMsg)
			}
		})
	}
}
