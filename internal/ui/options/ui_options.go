// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package options

type UiOptions struct {
	Plain       bool `cnfrmz:"plain"`
	BeautySleep bool `cnfrmz:"beauty-sleep"`
	Timestamps  bool `cnfrmz:"timestamps"`
}
