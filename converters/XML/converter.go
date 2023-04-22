package XML

import (
	"encoding/xml"
	"fmt"

	"serialization_tester/converters"
)

type Converter struct {
}

func (c *Converter) Serialize(p *converters.Person) ([]byte, error) {
	bytes, err := xml.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("xml marshal failed: %v", err)
	}
	return bytes, nil
}

func (c *Converter) Deserialize(raw []byte) (*converters.Person, error) {
	person := &converters.Person{}
	err := xml.Unmarshal(raw, person)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml bytes: %v", err)
	}
	return person, nil
}
