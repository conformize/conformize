// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package colors

type ColorCode string

const (
	reset  ColorCode = "0"
	red    ColorCode = "38;5;9"
	yellow ColorCode = "38;5;190"
	blue   ColorCode = "38;5;39"
	grey   ColorCode = "38;5;246"
)

type Color string

const (
	Red    Color = "red"
	Blue   Color = "blue"
	Yellow Color = "yellow"
	Grey   Color = "grey"
)

var colorMappings = map[Color]ColorCode{
	Red:    red,
	Blue:   blue,
	Yellow: yellow,
	Grey:   grey,
}

func (cn Color) String() string {
	return string(cn)
}

func (cn Color) Code() string {
	return string(colorMappings[cn])
}
