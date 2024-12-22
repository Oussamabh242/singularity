package encoder

import (
	"encoding/json"
	"testing"

	"github.com/Oussamabh242/singularity/pkg/parser"
)

func TestParse(t *testing.T) {
	t.Run("encoding something", func(t *testing.T) {
		packettype := 1
		metadata := parser.Metadata{
			Queue:       "thing",
			ContentType: "JSON",
			Topic:       "thing",
		}
		metadataStruct, _ := json.Marshal(metadata)
		message := ""
		res := Encode(uint8(packettype), metadataStruct, []byte(message))
    parse := parser.Parse(res)    

		t.Errorf("\n\n%v\n%v\n%v", res, string(metadataStruct),parse )

	})
}
