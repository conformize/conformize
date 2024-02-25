// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package benchmarks

import (
	"testing"

	"github.com/conformize/conformize/common/functions"
	"github.com/conformize/conformize/common/mocks"
)

func BenchmarkMergeTrees(b *testing.B) {
	printTrees := false
	trees := mocks.MockedConfigTrees(printTrees)
	for n := 0; n < b.N; n++ {
		functions.MergeTrees(trees...)
	}
}
