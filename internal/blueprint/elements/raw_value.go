// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package elements

type RawValue struct {
	Value     interface{}
	Sensitive bool
}

func (r *RawValue) GetValue() interface{} {
	return r.Value
}

func (r *RawValue) IsSensitive() bool {
	return r.Sensitive
}

func (r *RawValue) MarkSensitive() {
	r.Sensitive = true
}
