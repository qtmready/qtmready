// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

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
