// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package serialization

import "bytes"

type BufferedData struct {
	buf *bytes.Buffer
}

func (bdc *BufferedData) Read() ([]byte, error) {
	return bdc.buf.Bytes(), nil
}

func NewBufferedData(b []byte) SourceDataReader {
	return &BufferedData{
		buf: bytes.NewBuffer(b),
	}
}
