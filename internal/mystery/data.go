package mystery

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"fmt"
	"io/ioutil"
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
	return buf.Bytes(), nil
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	z := zlib.NewWriter(&b)
	if _, err := z.Write(data); err != nil {
		return nil, fmt.Errorf("unable to compress data: %v", err)
	}
	if err := z.Close(); err != nil {
		return nil, fmt.Errorf("unable to close compress: %v", err)
	}
	return b.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("unable to decompress data: %v", err)
	}
	result, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("decompress: unable to read data: %v", err)
	}
	return result, nil
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
