package yaml

import (
	"fmt"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"github.com/conformize/conformize/serialization/unmarshal/functions"
	"gopkg.in/yaml.v2"
)

type YamlUnmarshal struct{}

func (yamlUnmarshal *YamlUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, interface{}], error) {
	fileContent, err := source.Read()
	if err != nil {
		return nil, err
	}

	var content yaml.MapSlice
	err = yaml.Unmarshal(fileContent, &content)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml content: %v", err)
	}

	root := ds.NewNode[string, interface{}]()
	for _, item := range content {
		key := item.Key.(string)
		child := root.AddChild(key)
		unmarshalItem(item.Value, child)
	}
	return root, nil
}

func unmarshalItem(value interface{}, node *ds.Node[string, interface{}]) {
	switch val := value.(type) {
	case yaml.MapSlice:
		for _, item := range val {
			key := item.Key.(string)
			child := node.AddChild(key)
			unmarshalItem(item.Value, child)
		}
	case []interface{}:
		for _, elem := range val {
			if innerMap, ok := elem.(yaml.MapSlice); ok {
				for _, innerItem := range innerMap {
					key := innerItem.Key.(string)
					child := node.AddChild(key)
					unmarshalItem(innerItem.Value, child)
				}
			} else {
				functions.UnmarshalValue(node, val)
				break
			}
		}
	default:
		functions.UnmarshalValue(node, val)
	}
}
