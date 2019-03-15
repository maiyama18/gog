package git

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type Config struct {
	RepositoryFormatVersion int
	FileMode                bool
	Bare                    bool
}

func NewConfig(confPath string) (*Config, error) {
	f, err := os.Open(confPath)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(f)

	// default conf
	conf := &Config{
		RepositoryFormatVersion: 0,
		FileMode:                true,
		Bare:                    true,
	}

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, "repositoryformatversion") {
			elems := strings.Split(line, "=")
			if len(elems) == 2 {
				vs := elems[1]
				v, err := strconv.Atoi(vs)
				if err == nil {
					conf.RepositoryFormatVersion = v
				}
			}
		}
		if strings.HasPrefix(line, "filemode") {
			elems := strings.Split(line, "=")
			if len(elems) == 2 {
				ms := elems[1]
				m, err := strconv.ParseBool(ms)
				if err == nil {
					conf.FileMode = m
				}
			}
		}
		if strings.HasPrefix(line, "bare") {
			elems := strings.Split(line, "=")
			if len(elems) == 2 {
				bs := elems[1]
				b, err := strconv.ParseBool(bs)
				if err == nil {
					conf.Bare = b
				}
			}
		}
	}

	return conf, nil
}

type Repository struct {
	WorkTree string
	GitDir   string
	Conf     *Config
}

// force disables all check for initializing repository, which is used when 'git init'.
func NewRepository(workTree string, force bool) (*Repository, error) {
	gitDir := path.Join(workTree, ".git")

	if !force && !isExistingDir(gitDir) {
		return nil, fmt.Errorf("not a git repository: %s", gitDir)
	}

	repo := &Repository{WorkTree: workTree, GitDir: gitDir}

	confPath, err := repo.gitFile(false, "config")
	if err != nil {
		return nil, err
	}

	if err := repo.readConf(force, confPath); err != nil {
		return nil, err
	}

	return repo, nil
}

// file returns
func (r *Repository) gitFile(mkdir bool, elems ...string) (string, error) {
	if len(elems) == 0 {
		return "", errors.New("filepath not provided")
	}

	elems = append([]string{r.GitDir}, elems[:len(elems)-1]...)
	dirPath, err := r.gitDir(mkdir, elems...)
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

func (r *Repository) gitDir(mkdir bool, elems ...string) (string, error) {
	elems = append([]string{r.GitDir}, elems...)
	dirPath := path.Join(elems...)

	if !isExistingDir(dirPath) {
		if !mkdir {
			return "", fmt.Errorf("not a directory: %s", dirPath)
		}
		if err := os.MkdirAll(dirPath, 0777); err != nil {
			return "", err
		}
		return dirPath, nil
	}
	return dirPath, nil
}

func (r *Repository) readConf(force bool, confPath string) error {
	if !isExistingFile(confPath) {
		if !force {
			return fmt.Errorf("config file not found: %s", confPath)
		}
		return nil
	}
	conf, err := NewConfig(confPath)
	if err != nil {
		return err
	}
	r.Conf = conf
	return nil
}

func isExistingDir(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		// file does not exist
		return false
	}

	return fi.Mode().IsDir()
}

func isExistingFile(filePath string) bool {
	fi, err := os.Stat(filePath)
	if err != nil {
		// file does not exist
		return false
	}

	return fi.Mode().IsRegular()
}
