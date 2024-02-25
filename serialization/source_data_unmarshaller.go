package serialization

import (
	"github.com/conformize/conformize/common/ds"
)

type SourceDataUnmarshaller interface {
	Unmarshal(srcDataRdr SourceDataReader) (*ds.Node[string, interface{}], error)
}
