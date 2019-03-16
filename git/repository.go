package git

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

const defaultDescription = "Unnamed repository; edit this file 'description' to name the repository.\n"

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

	if isExistingFile(workTree) {
		return nil, fmt.Errorf("work tree already exist as regular file: %s", workTree)
	}

	if isExistingDir(workTree) {
		files, err := ioutil.ReadDir(workTree)
		if err != nil {
			return nil, err
		}

		// fail if workTree has files other than .git
		if len(files) > 1 {
			fileNames := make([]string, 0)
			for _, file := range files {
				fileNames = append(fileNames, file.Name())
			}
			return nil, fmt.Errorf("work tree is not empty: %s (%v)", workTree, fileNames)
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
	if err := repo.writeToGitFile(repo.Conf.format(), "config"); err != nil {
		return nil, err
	}

	return repo, nil
}

func FindRepository(curPath string, required bool) (*Repository, error) {
	if isExistingDir(path.Join(curPath, ".git")) {
		return NewRepository(curPath, false)
	}

	if curPath == "/" {
		if required {
			return nil, errors.New("couldn't find git repository")
		}
		return nil, nil
	}

	return FindRepository(path.Dir(curPath), required)
}

// force disables all check for initializing repository, which is used when 'git init'.
func NewRepository(workTree string, force bool) (*Repository, error) {
	gitDir := path.Join(workTree, ".git")

	if !force && !isExistingDir(gitDir) {
		return nil, fmt.Errorf("not a git repository: %s", gitDir)
	}

	repo := &Repository{WorkTree: workTree, GitDir: gitDir, Conf: &defaultConfig}

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
	f, err := os.Create(filePath)
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
	conf, err := newConfig(confFile)
	if err != nil {
		return err
	}
	r.Conf = conf
	return nil
}

func (r *Repository) ReadObject(sha string, expectedKind string) (Object, error) {
	objPath, err := r.gitFile(false, "objects", sha[:2], sha[2:])
	if err != nil {
		return nil, err
	}

	f, err := os.Open(objPath)
	if err != nil {
		return nil, err
	}
	rd, err := zlib.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	raw, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}

	si := bytes.Index(raw, []byte(" "))
	if si < 0 {
		return nil, fmt.Errorf("contents of %s does not contain space: %v", sha, raw)
	}
	kind := string(raw[0:si])
	if kind != expectedKind {
		return nil, fmt.Errorf("object type mismatch: provided=%s, got=%s", expectedKind, kind)
	}
	raw = raw[si+1:]

	ni := bytes.Index(raw, []byte{0})
	if ni < 0 {
		return nil, fmt.Errorf("contents of %s does not contain null char: %v", sha, raw)
	}
	size, err := strconv.Atoi(string(raw[0:ni]))
	if err != nil {
		return nil, err
	}
	raw = raw[ni+1:]
	if size != len(raw) {
		return nil, fmt.Errorf("wrong length of object %s: header says=%d, actual=%d", sha, size, len(raw))
	}

	switch expectedKind {
	case "commit":
	case "tree":
	case "tag":
	case "blob":
		return NewBlob(string(raw)), nil
	}
	return nil, fmt.Errorf("unknown kind for object %v: %s", sha, kind)
}

func (r *Repository) WriteObject(obj Object, dryRun bool) (string, error) {
	data := obj.Serialize()

	hw := sha1.New()
	_, _ = io.WriteString(hw, obj.Kind())
	_, _ = io.WriteString(hw, " ")
	_, _ = io.WriteString(hw, strconv.Itoa(len(data)))
	hw.Write([]byte{0})
	_, _ = io.WriteString(hw, data)

	sha := fmt.Sprintf("%x", hw.Sum(nil))

	if !dryRun {
		if err := r.writeToGitFile(sha, "objects", sha[0:2], sha[2:]); err != nil {
			return "", err
		}
	}

	return sha, nil
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
