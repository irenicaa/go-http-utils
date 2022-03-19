package httputils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// GetJSONData ...
func GetJSONData(reader io.Reader, data interface{}) error {
	dataAsJSON, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("unable to read the JSON data: %s", err)
	}

	if err := json.Unmarshal(dataAsJSON, data); err != nil {
		return fmt.Errorf("unable to unmarshal the JSON data: %s", err)
	}

	return nil
}
