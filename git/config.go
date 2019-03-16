package git

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Config struct {
	RepositoryFormatVersion int
	FileMode                bool
	Bare                    bool
}

var DefaultConfig = Config{
	RepositoryFormatVersion: 0,
	FileMode:                true,
	Bare:                    false,
}

func NewConfig(confFile io.Reader) (*Config, error) {
	sc := bufio.NewScanner(confFile)

	// default conf
	conf := &DefaultConfig
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "repositoryformatversion") {
			elems := strings.Split(line, "=")
			if len(elems) == 2 {
				vs := strings.TrimSpace(elems[1])
				v, err := strconv.Atoi(vs)
				if err == nil {
					conf.RepositoryFormatVersion = v
				}
			}
		}
		if strings.HasPrefix(line, "filemode") {
			elems := strings.Split(line, "=")
			if len(elems) == 2 {
				ms := strings.TrimSpace(elems[1])
				m, err := strconv.ParseBool(ms)
				if err == nil {
					conf.FileMode = m
				}
			}
		}
		if strings.HasPrefix(line, "bare") {
			elems := strings.Split(line, "=")
			if len(elems) == 2 {
				bs := strings.TrimSpace(elems[1])
				b, err := strconv.ParseBool(bs)
				if err == nil {
					conf.Bare = b
				}
			}
		}
	}

	return conf, nil
}

func (c *Config) Format() string {
	return fmt.Sprintf(`[core]
	repositoryformatversion = %d
	filemode = %t
	bare = %t`, c.RepositoryFormatVersion, c.FileMode, c.Bare)
}
