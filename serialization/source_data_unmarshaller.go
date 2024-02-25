package serialization

import (
	"github.com/conformize/conformize/common/ds"
)

type SourceDataUnmarshaller interface {
	Unmarshal(srcDataRdr SourceDataReader) (*ds.Node[string, any], error)
}
