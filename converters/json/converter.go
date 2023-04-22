package json

import (
	"encoding/json"
	"fmt"

	"serialization_tester/converters"
)

type Converter struct {
}

func (c *Converter) Serialize(p *converters.Person) ([]byte, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed: %v", err)
	}
	return bytes, nil
}

func (c *Converter) Deserialize(raw []byte) (*converters.Person, error) {
	person := &converters.Person{}
	err := json.Unmarshal(raw, person)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json bytes: %v", err)
	}
	return person, nil
}
