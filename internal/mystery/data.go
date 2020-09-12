package mystery

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Data defines object for store
type Data map[string]interface{}

// Encode provides encoding of data
func Encode(d Data) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode provides decoding of data into Data representation
func Decode(data []byte) (Data, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var out Data
	if err := dec.Decode(&out); err != nil {
		return out, fmt.Errorf("unable to decode value: %v", err)
	}
	return out, nil
}
