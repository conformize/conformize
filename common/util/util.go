// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package util

import (
	"os"
	"path/filepath"
	"strings"
)

func ResolveFileBasePath(filePath string) (string, error) {
	resolvedPath, err := ResolveFilePath(filePath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(resolvedPath), nil
}

func ResolveFilePath(filePath string) (string, error) {
	baseDir := GetWorkDir()
	return ResolveFileRelativePath(baseDir, filePath)
}

func isSymlink(path string) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}

func ResolveFileRelativePath(baseDir, filePath string) (string, error) {
	var err error
	var resolvedPath = filePath
	if strings.HasPrefix(filePath, "~/") || filePath == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		resolvedPath = filepath.Join(home, strings.TrimPrefix(filePath, "~"))
	}

	if !filepath.IsAbs(resolvedPath) {
		resolvedPath = filepath.Join(baseDir, resolvedPath)
	}

	var pathIsSymlink bool
	resolvedPath = filepath.Clean(resolvedPath)
	if pathIsSymlink, err = isSymlink(resolvedPath); pathIsSymlink && err == nil {
		resolvedPath, err = filepath.EvalSymlinks(resolvedPath)
		if err != nil {
			return "", err
		}
	}

	if _, err = os.Stat(resolvedPath); err != nil {
		return "", err
	}

	return resolvedPath, nil
}

var workDir string

func SetWorkDir(path string) {
	workDir = path
}

func GetWorkDir() string {
	if len(workDir) == 0 {
		var err error
		workDir, err = os.Getwd()
		if err != nil {
			workDir = "."
		}
	}
	return workDir
}
