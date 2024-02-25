// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package serialization

import (
	"os"

	"github.com/conformize/conformize/common/util"
)

type FileSource struct {
	FilePath string
}

func (c *FileSource) Read() ([]byte, error) {
	return os.ReadFile(c.FilePath)
}

func NewFileSource(filePath string) (*FileSource, error) {
	resolvedPath, err := util.ResolveFilePath(filePath)
	if err != nil {
		return nil, err
	}
	return &FileSource{FilePath: resolvedPath}, nil
}
