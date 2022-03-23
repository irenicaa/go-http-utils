package httputils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// HTTPClient ...
type HTTPClient interface {
	Do(request *http.Request) (*http.Response, error)
}

// ReadJSONData ...
func ReadJSONData(reader io.Reader, data interface{}) error {
	dataAsJSON, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("unable to read the JSON data: %w", err)
	}

	if err := json.Unmarshal(dataAsJSON, data); err != nil {
		return fmt.Errorf("unable to unmarshal the JSON data: %w", err)
	}

	return nil
}

// LoadJSONData ...
func LoadJSONData(
	httpClient HTTPClient,
	url string,
	authHeader string,
	responseData interface{},
) error {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("unable to create the request: %w", err)
	}

	if authHeader != "" {
		request.Header.Add("Authorization", authHeader)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("unable to send the request: %w", err)
	}
	defer response.Body.Close()

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("unable to read the request body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"request was failed: %d %s",
			response.StatusCode,
			responseBytes,
		)
	}

	if err = json.Unmarshal(responseBytes, responseData); err != nil {
		return fmt.Errorf("unable to unmarshal the request body: %w", err)
	}

	return nil
}
