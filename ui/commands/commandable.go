// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package commands

import "flag"

type Commandable interface {
	GetExpression() string
	GetDescription() string
	IsHidden() bool
	GetFlags() *flag.FlagSet
	GetSubcommands() []CommandRunnable
	GetAliases() []string
}
