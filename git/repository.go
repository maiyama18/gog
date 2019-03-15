package git

import (
	"errors"
	"fmt"
	"os"
	"path"
)

type Config struct {
}

type Repository struct {
	WorkTree string
	GitDir   string
	Conf     Config
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
	// TODO: configファイルから設定を読み込む
	r.Conf = Config{}
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
