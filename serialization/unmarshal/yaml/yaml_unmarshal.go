package yaml

import (
	"fmt"

	"github.com/conformize/conformize/common/ds"
	"github.com/conformize/conformize/serialization"
	"gopkg.in/yaml.v2"
)

type YamlUnmarshal struct{}

func (yamlUnmarshal *YamlUnmarshal) Unmarshal(source serialization.SourceDataReader) (*ds.Node[string, any], error) {
	fileContent, err := source.Read()
	if err != nil {
		return nil, err
	}

	var root = ds.NewNode[string, any]()
	err = yaml.Unmarshal(fileContent, &root)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml content: %v", err)
	}
	return root, nil
}
