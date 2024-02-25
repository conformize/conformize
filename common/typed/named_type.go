// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package typed

type NamedType string

const (
	BooleanType NamedType = "boolean"
	NumberType  NamedType = "number"
	StringType  NamedType = "string"
	ListType    NamedType = "list"
	MapType     NamedType = "map"
	ObjectType  NamedType = "object"
	TupleType   NamedType = "tuple"
	VariantType NamedType = "variant"
	GenericType NamedType = "generic"
	InvalidType NamedType = "invalid"
)

func NamedTypeFromTypeHint(hint TypeHint) string {
	switch hint {
	case Boolean:
		return string(BooleanType)
	case Number:
		return string(NumberType)
	case String:
		return string(StringType)
	case List:
		return string(ListType)
	case Map:
		return string(MapType)
	case Object:
		return string(ObjectType)
	case Tuple:
		return string(TupleType)
	case Variant:
		return string(VariantType)
	case Generic:
		return string(GenericType)
	default:
		return string(InvalidType)
	}
}
