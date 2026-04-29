package internal

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// MarshalBinaryValue serializes a value to binary bytes using gob.
func MarshalBinaryValue(value any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(value); err != nil {
		return nil, fmt.Errorf("marshal binary value: %w", err)
	}
	return buffer.Bytes(), nil
}

// UnmarshalBinaryValue deserializes gob bytes into value.
func UnmarshalBinaryValue[T any](data []byte, value *T) error {
	if value == nil {
		return fmt.Errorf("unmarshal binary value: nil target")
	}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("unmarshal binary value: %w", err)
	}
	return nil
}
