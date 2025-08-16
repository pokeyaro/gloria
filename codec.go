package gloria

import (
	"encoding/json"
)

// Codec defines how to marshal/unmarshal request/response payloads.
type Codec interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

// stdJSONCodec is the default codec using encoding/json.
type stdJSONCodec struct{}

func (stdJSONCodec) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (stdJSONCodec) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
