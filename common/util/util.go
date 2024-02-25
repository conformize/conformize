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
	completeFilePath, err := ResolveFilePath(filePath)
	if err != nil {
		return "", err
	}
	return filepath.Dir(completeFilePath), nil
}

func ResolveFilePath(filePath string) (string, error) {
	var err error
	var basePath string

	path := filePath
	basePath, err = os.Getwd()
	if err != nil {
		return "", err
	}

	if !filepath.IsAbs(path) {
		if strings.HasPrefix(path, "~/") {
			if basePath, err = os.UserHomeDir(); err == nil {
				path = path[2:]
			}
		}

		if err != nil {
			return "", err
		}
		path = filepath.Join(basePath, path)
	}

	path = filepath.ToSlash(path)
	path, err = filepath.EvalSymlinks(path)
	if _, err = os.Stat(path); err != nil {
		return "", err
	}
	return path, err
}
