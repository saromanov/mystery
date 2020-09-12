package mystery

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"fmt"
)

// Data defines object for store
type Data string

// encode provides encoding of data
func (d Data) encode() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(d)
	if err != nil {
		return nil, err
	}
	return buf.Bytes()
}

func comperss(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := zlib.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, fmt.Errorf("unable to compress data: %v", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("unable to close compress: %v", err)
	}
	fmt.Println("LENAFTER: ", len(b.Bytes()))
	return b.Bytes(), nil
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
