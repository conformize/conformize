// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package colors

type colorCode string

const (
	resetCode  colorCode = "0"
	redCode    colorCode = "38;5;9"
	yellowCode colorCode = "38;5;190"
	blueCode   colorCode = "38;5;39"
	greyCode   colorCode = "38;5;246"
	greenCode  colorCode = "38;5;34"
)

type Color int

const (
	Red Color = iota
	Blue
	Yellow
	Grey
	Green
)

var colorNames = []string{
	"red",
	"blue",
	"yellow",
	"grey",
	"green",
}

var colorCodes = []colorCode{
	redCode,
	blueCode,
	yellowCode,
	greyCode,
	greenCode,
}

func (c Color) String() string {
	if int(c) < 0 || int(c) >= len(colorNames) {
		return "unknown"
	}
	return colorNames[c]
}

func (c Color) Code() string {
	if int(c) < 0 || int(c) >= len(colorCodes) {
		return string(resetCode)
	}
	return string(colorCodes[c])
}
