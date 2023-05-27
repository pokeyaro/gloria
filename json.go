package gloria

import (
	"encoding/json"

	gojson "github.com/goccy/go-json"
)

// JSONLibrary Define the interface for serialization and deserialization of the json parsing library.
type JSONLibrary interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// NativeJSONLibrary is the native implementation of encoding/json.
type NativeJSONLibrary struct{}

func (l NativeJSONLibrary) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (l NativeJSONLibrary) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// GoJSONLibrary is an implementation of the popular tripartite library go-json.
type GoJSONLibrary struct{}

func (l GoJSONLibrary) Marshal(v interface{}) ([]byte, error) {
	return gojson.Marshal(v)
}

func (l GoJSONLibrary) Unmarshal(data []byte, v interface{}) error {
	return gojson.Unmarshal(data, v)
}

// json-iterator implementation
//type JSONIteratorLibrary struct{}
//
//func (l JSONIteratorLibrary) Marshal(v interface{}) ([]byte, error) {
//	return jsoniter.Marshal(v)
//}
//
//func (l JSONIteratorLibrary) Unmarshal(data []byte, v interface{}) error {
//	return jsoniter.Unmarshal(data, v)
//}

// bytedance/sonic implementation
//type SonicLibrary struct{}
//
//func (l SonicLibrary) Marshal(v interface{}) ([]byte, error) {
//	return sonic.Marshal(v)
//}
//
//func (l SonicLibrary) Unmarshal(data []byte, v interface{}) error {
//	return sonic.Unmarshal(data, v)
//}
