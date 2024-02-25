package blueprint

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func marshalToYAML(v any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)

	encoder.SetIndent(2)
	err := encoder.Encode(v)
	if err != nil {
		return []byte{}, err
	}

	encoder.Close()
	return buf.Bytes(), err
}

func marshalToJSON(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func NewBlueprintMarshaller(fileName string, blprnt *Blueprint) *BlueprintMarshaller {
	var marshalFn func(v any) ([]byte, error)
	if strings.HasSuffix(fileName, "cnfrm.yaml") {
		marshalFn = marshalToYAML
	}

	if strings.HasSuffix(fileName, "cnfrm.json") {
		marshalFn = marshalToJSON
	}

	return &BlueprintMarshaller{
		fileName:  fileName,
		blueprint: blprnt,
		marshalFn: marshalFn,
	}
}

type BlueprintMarshaller struct {
	fileName  string
	blueprint *Blueprint
	marshalFn func(v any) ([]byte, error)
}

func (blprntMarshaller *BlueprintMarshaller) Marshal() error {
	file, fileErr := os.Create(blprntMarshaller.fileName)
	if fileErr != nil {
		return fileErr
	}

	defer file.Close()
	data, err := blprntMarshaller.marshalFn(blprntMarshaller.blueprint)
	if err != nil {
		return err
	}

	if _, err = file.Write(data); err != nil {
		return err
	}
	return nil
}
