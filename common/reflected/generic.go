package reflected

import (
	"reflect"

	"github.com/conformize/conformize/common/typed"
)

func Generic(val reflect.Value) (typed.Valuable, error) {
	valTypeHint := typed.TypeHintOf(val)
	value, err := ValueFromTypeHint(val, valTypeHint)
	if err != nil {
		return nil, err
	}
	return &typed.GenericValue{Value: value}, nil
}
