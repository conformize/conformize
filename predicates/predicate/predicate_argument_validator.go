package predicate

import (
	"fmt"

	"github.com/conformize/conformize/common/typed"
	"github.com/conformize/conformize/internal/providers/api/schema"
	"github.com/conformize/conformize/internal/providers/api/schema/attributes"
)

type PredicateArgumentsValidator struct{}

func (prdArgValidator *PredicateArgumentsValidator) Validate(value typed.Valuable, args *typed.TupleValue, schema schema.Schemable) (bool, error) {
	prdAttrs := schema.GetAttributes()
	if prdValAttr, ok := prdAttrs["Value"]; ok {
		prdValAttrTypeHint := prdValAttr.Type().Hint()
		if prdValAttrTypeHint != typed.Generic && prdValAttrTypeHint != value.Type().Hint() {
			return false, fmt.Errorf(
				"wrong value type - expected type %s, got %s", prdValAttr.Type().Name(), value.Type().Name(),
			)
		}
	}

	if prdArgsAttr, ok := prdAttrs["Arguments"]; ok {
		schemaArgs, ok := prdArgsAttr.(*attributes.TupleAttribute)
		if !ok {
			return false, fmt.Errorf(
				"wrong arguments type specified in schema - expected %s, got %s", typed.TupleType, prdArgsAttr.Type().Name(),
			)
		}

		argsLen := len(schemaArgs.GetElementsTypes())
		if argsLen > 0 {
			if args == nil || len(args.Elements) == 0 {
				return false, fmt.Errorf("expected %d arguments, got %v", argsLen, args)
			}

			providedArgsLen := len(args.Elements)
			if argsLen > providedArgsLen {
				return false, fmt.Errorf("wrong number of arguments provided - expected %d arguments, got %d", argsLen, providedArgsLen)
			}

			for argIdx, schemaArgType := range schemaArgs.GetElementsTypes() {
				if schemaArgType.Hint() == typed.Generic {
					continue
				}

				providedArg := args.Elements[argIdx]
				if schemaArgType.Hint() != providedArg.Type().Hint() {
					return false,
						fmt.Errorf("wrong type for argument [%d] - %s expected, got %s",
							argIdx, schemaArgs.Type().Name(), providedArg.Type().Name(),
						)
				}
			}
		}
	}

	return true, nil
}
