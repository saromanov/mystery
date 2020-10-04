package mystery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/saromanov/mystery/internal/backend"
)

var (
	errNoMasterPass = errors.New("master pass is not defined")
	errNoNamespace  = errors.New("namespace is not defined")
	errNoValue      = errors.New("value is not defined")
	errNoBackend    = errors.New("backend is not defined")
)

// PutRequest provides struct for sending to Put
type PutRequest struct {
	MasterPass string
	Namespace  string
	Data       Data
	Type       string
	Backend    backend.Backend
}

// validate provides validation of data
func (p PutRequest) validate() error {
	if p.MasterPass == "" {
		return errNoMasterPass
	}
	if p.Namespace == "" {
		return errNoNamespace
	}
	if len(p.Data) == 0 && p.Type != "file" {
		return errNoValue
	}
	if p.Backend == nil {
		return errNoBackend
	}
	return nil
}

// Put provides adding key-value pair to backend
func (m *Mystery) Put(p PutRequest) error {
	if err := p.validate(); err != nil {
		return fmt.Errorf("put: unable to validate data: %v", err)
	}

	if p.Type == "file" {
		fmt.Println("DATA: ", string(p.Data))
		data, err := readFile(string(p.Data))
		if err != nil {
			return fmt.Errorf("unable to read file: %v", err)
		}
		p.Data = Data(data)
	}

	value, err := p.Data.encode()
	if err != nil {
		return fmt.Errorf("unable to encode value: %v", err)
	}
	reducedValue, err := compress(value)
	if err != nil {
		return fmt.Errorf("unable to compress data: %v", err)
	}
	compressed := len(reducedValue) < len(value)
	if compressed {
		value = reducedValue
	}
	if err := p.Backend.Put([]byte(p.MasterPass), backend.Secret{
		Namespace:  []byte(p.Namespace),
		Data:       value,
		Compressed: compressed,
	}); err != nil {
		return fmt.Errorf("put: unable to store data: %v", err)
	}
	return nil
}

// readFile provides reading of the json file and convert it to format key=value
func readFile(name string) (string, error) {
	data, err := unmarshal(name)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal data: %v", err)
	}
	result := ""
	for k, v := range data {
		result += fmt.Sprintf("%s=%s;", k, v.(string))
	}
	return result, nil
}

func unmarshal(name string) (map[string]interface{}, error) {
	jsonFile, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s %v", name, err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(byteValue), &result); err != nil {
		return nil, fmt.Errorf("unable to unmarshal data: ", err)
	}
	return result, nil
}
