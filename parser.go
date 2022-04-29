package github_webhook

import (
	"encoding/json"
	"fmt"
	"github.com/luciddev13/limit_reader"
	"io/ioutil"
	"net/http"
)

// Loads the raw payload from the request.
//
// Since we need to load the entire body in order to unmarshal it
// this will limit the read to a size that can be reasonably expected
// to be sent by GitHub, anything larger is truncated.
func loadRawPayload(r *http.Request, maxSize int) ([]byte, error) {
	limitReader := limit_reader.New(r.Body, int64(maxSize))
	rawPayload, err := ioutil.ReadAll(limitReader)
	if err != nil {
		return nil, fmt.Errorf("could not read payload: %v", err)
	}

	return rawPayload, nil
}

// Parses the fully validated raw payload
func parsePayload(rawPayload []byte) (map[string]interface{}, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		return nil, fmt.Errorf("could not unmarshal payload: %v", err)
	}
	return payload, nil
}
