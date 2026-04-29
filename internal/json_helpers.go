package internal

import (
	"encoding/json"
	"fmt"
)

// MarshalJSONValue serializes a value to JSON bytes.
func MarshalJSONValue(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal json value: %w", err)
	}
	return data, nil
}
