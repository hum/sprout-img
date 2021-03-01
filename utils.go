package sproutimg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status ResponseType
	Data   map[string]string
}

type ResponseType string

const (
	DEFAULT ResponseType = "DEFAULT"
	OK      ResponseType = "OK"
	ERROR   ResponseType = "ERROR"
)

func setHeader(w http.ResponseWriter, values map[string]string) error {
	if len(values) == 0 {
		return fmt.Errorf("setHeader() could not perform any action. Parameter 'values' is empty.")
	}
	for k, v := range values {
		w.Header().Set(k, v)
	}
	return nil
}

func marshalResponse(data interface{}) ([]byte, error) {
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json, nil
}
