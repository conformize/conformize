// Copyright (c) 2024 Hristo Paskalev
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package hcl

import (
	"fmt"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type HclFileUnmarshal struct{}

func (hclUnmarshal *HclFileUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, any], error) {
	hclParser := hclparse.NewParser()
	content, err := source.Read()
	if err != nil {
		return nil, err
	}
	file, diags := hclParser.ParseHCL(content, "hcl_file")
	if diags.HasErrors() {
		return nil, fmt.Errorf("HCL parsing error: %s", diags.Error())
	}

	root := ds.NewNode[string, any]()
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("expected HCL body, got %T", file.Body)
	}

	err = parseBody(body, root)
	return root, err
}

func parseBody(body *hclsyntax.Body, parent *ds.Node[string, any]) error {
	emptyCtx := &hcl.EvalContext{}
	for name, attr := range body.Attributes {
		val, diag := attr.Expr.Value(emptyCtx)
		if diag.HasErrors() {
			return fmt.Errorf("error reading attribute %s: %s", name, diag.Error())
		}
		decodeValue(val, parent, name)
	}

	for _, block := range body.Blocks {
		current := parent

		current = current.AddChild(block.Type)

		for _, label := range block.Labels {
			current = current.AddChild(label)
		}

		if err := parseBody(block.Body, current); err != nil {
			return err
		}
	}
	return nil
}

func decodeValue(val cty.Value, parent *ds.Node[string, any], attrName string) {
	if !val.IsKnown() || val.IsNull() {
		if attrName != "" {
			parent.AddAttribute(attrName, nil)
		}
		return
	}

	switch {
	case val.Type().IsPrimitiveType():
		var v any
		switch val.Type() {
		case cty.String:
			v = val.AsString()
		case cty.Number:
			f, _ := val.AsBigFloat().Float64()
			v = f
		case cty.Bool:
			v = val.True()
		}
		if attrName != "" {
			parent.AddAttribute(attrName, v)
		}

	case val.Type().IsMapType() || val.Type().IsObjectType():
		var objNode *ds.Node[string, any]
		if attrName != "" {
			objNode = parent.AddChild(attrName)
		} else {
			objNode = parent
		}
		for k, v := range val.AsValueMap() {
			decodeValue(v, objNode, k)
		}
	case val.Type().IsListType() || val.Type().IsTupleType() || val.Type().IsSetType():
		if attrName == "" {
			return
		}
		values := []any{}
		for _, item := range val.AsValueSlice() {
			switch {
			case item.Type().IsPrimitiveType():
				switch item.Type() {
				case cty.String:
					values = append(values, item.AsString())
				case cty.Number:
					f, _ := item.AsBigFloat().Float64()
					values = append(values, f)
				case cty.Bool:
					values = append(values, item.True())
				default:
					values = append(values, item.GoString())
				}
			default:
				child := ds.NewNode[string, any]()
				decodeValue(item, child, "")
				values = append(values, child.Value)
			}
		}
		parent.AddAttribute(attrName, values)

	default:
		if attrName != "" {
			parent.AddAttribute(attrName, val.GoString())
		}
	}
}
