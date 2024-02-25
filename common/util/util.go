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

func FilePath(filePath string) (string, error) {
	path := filepath.ToSlash(filePath)

	if !filepath.IsAbs(path) {
		var basePath string
		var err error

		if !strings.HasPrefix(path, "~/") {
			basePath, err = os.Getwd()
		} else if basePath, err = os.UserHomeDir(); err == nil {
			path = path[2:]
		}

		if err != nil {
			return "", err
		}
		path = filepath.Join(basePath, path)
	}

	path, err := filepath.EvalSymlinks(path)
	return path, err
}
