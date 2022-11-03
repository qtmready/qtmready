package utils

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrNotGitRepo = errors.New("ctrlplane must be initialized in a git repository")
	ErrUnknown    = errors.New("unknown error")
)

func DetectGitRoot() (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// recursively go to the until we find the git root or throw an error
	for {
		info, err := os.Stat(path.Join(workdir, ".git"))
		if err == nil {
			if !info.IsDir() {
				return "", ErrNotGitRepo
			}

			return path.Join(workdir, ".git"), nil
		}

		if !os.IsNotExist(err) {
			// unknown error
			return "", err
		}

		// detect bare repo
		ok, err := isGitDir(workdir)
		if err != nil {
			return "", err
		}

		if ok {
			return workdir, nil
		}

		if parent := filepath.Dir(workdir); parent == workdir {
			return "", fmt.Errorf(".git not found")
		}
	}
}

func isGitDir(cwd string) (bool, error) {
	markers := []string{"HEAD", "objects", "refs"}

	for _, marker := range markers {
		_, err := os.Stat(path.Join(cwd, marker))
		if err == nil {
			continue
		}

		if !os.IsNotExist(err) {
			return false, ErrUnknown
		}

		return false, nil
	}

	return true, nil
}
