// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package functions

import commonfns "github.com/conformize/conformize/common/functions"

func EnvVarLookup(varName string) func() any {
	return func() any {
		envVar, _ := commonfns.LookupEnvVar(varName)
		return envVar
	}
}
