package interval

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
)

func marshalJSONValue(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal json value: %w", err)
	}
	return data, nil
}

func jsonResultString(data []byte, err error, fallback string) string {
	if err != nil {
		return fallback
	}
	return string(data)
}

func marshalBinaryValue(value any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(value); err != nil {
		return nil, fmt.Errorf("marshal binary value: %w", err)
	}
	return buffer.Bytes(), nil
}

func unmarshalBinaryValue(data []byte, value any) error {
	if value == nil {
		return fmt.Errorf("unmarshal binary value: nil target")
	}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(value); err != nil {
		return fmt.Errorf("unmarshal binary value: %w", err)
	}
	return nil
}
