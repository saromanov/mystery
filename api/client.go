package api

import (
	"fmt"
    "net/http"
    "bytes"
	"encoding/json"
)

// Post provides sending of the post request
func Post(url string, data interface{}) error {
	marshaled, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal data: %v", err)
	}
    resp, err := http.Post(url, "application/json",
        bytes.NewBuffer(marshaled))

    if err != nil {
        return fmt.Errorf("unable to send post request: %v", err)
    }
    var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

    return nil
}