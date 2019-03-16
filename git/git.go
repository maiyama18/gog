package git

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const defaultDescription = "Unnamed repository; edit this file 'description' to name the repository.\n"

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

type Repository struct {
	WorkTree string
	GitDir   string
	Conf     *Config
}

func CreateRepository(workTree string) (*Repository, error) {
	repo, err := NewRepository(workTree, true)
	if err != nil {
		return nil, err
	}

	// make sure that workTree is nonexistent or empty
	if isExistingFile(workTree) {
		return nil, fmt.Errorf("work tree already exist as regular file: %s", workTree)
	}

	if isExistingDir(workTree) {
		files, err := ioutil.ReadDir(workTree)
		if err != nil {
			return nil, err
		}
		if len(files) > 0 {
			return nil, fmt.Errorf("work tree is not empty: %s", workTree)
		}
	} else {
		if err := os.MkdirAll(workTree, 0777); err != nil {
			return nil, err
		}
	}

	// create subdirectories of .git if not exists
	if _, err := repo.gitDir(true, "branches"); err != nil {
		return nil, err
	}
	if _, err := repo.gitDir(true, "objects"); err != nil {
		return nil, err
	}
	if _, err := repo.gitDir(true, "refs", "tags"); err != nil {
		return nil, err
	}
	if _, err := repo.gitDir(true, "refs", "heads"); err != nil {
		return nil, err
	}

	if err := repo.writeToGitFile(defaultDescription, "description"); err != nil {
		return nil, err
	}
	if err := repo.writeToGitFile("ref: refs/heads/master\n", "HEAD"); err != nil {
		return nil, err
	}
	if err := repo.writeToGitFile(repo.Conf.Format(), "config"); err != nil {
		return nil, err
	}

	return repo, nil
}

// TODO: privateにする
// force disables all check for initializing repository, which is used when 'git init'.
func NewRepository(workTree string, force bool) (*Repository, error) {
	gitDir := path.Join(workTree, ".git")

	if !force && !isExistingDir(gitDir) {
		return nil, fmt.Errorf("not a git repository: %s", gitDir)
	}

	repo := &Repository{WorkTree: workTree, GitDir: gitDir, Conf: &DefaultConfig}

	confPath, err := repo.gitFile(false, "config")
	if !force && err != nil {
		return nil, err
	}

	if err := repo.readConf(force, confPath); !force && err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *Repository) gitFile(mkdir bool, relPath ...string) (string, error) {
	if len(relPath) == 0 {
		return "", errors.New("filepath not provided")
	}

	dirRelPath := relPath[:len(relPath)-1]
	dirPath, err := r.gitDir(mkdir, dirRelPath...)
	if err != nil {
		return "", err
	}
	return path.Join(dirPath, relPath[len(relPath)-1]), nil
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

func (r *Repository) writeToGitFile(content string, relPath ...string) error {
	filePath, err := r.gitFile(false, relPath...)
	if err != nil {
		return err
	}
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

func (r *Repository) readConf(force bool, confPath string) error {
	if !isExistingFile(confPath) {
		if !force {
			return fmt.Errorf("config file not found: %s", confPath)
		}
		return nil
	}

	confFile, err := os.Open(confPath)
	if err != nil {
		return err
	}
	conf, err := NewConfig(confFile)
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
