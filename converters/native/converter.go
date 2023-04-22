package native

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"serialization_tester/converters"
)

type Converter struct {
}

func (c *Converter) Serialize(p *converters.Person) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p); err != nil {
		return nil, fmt.Errorf("xml marshal failed: %v", err)
	}
	return buf.Bytes(), nil
}

func (c *Converter) Deserialize(raw []byte) (*converters.Person, error) {
	var buf bytes.Buffer
	buf.Write(raw)

	person := &converters.Person{}
	dec := gob.NewDecoder(&buf)
	if err := dec.Decode(person); err != nil {
		return nil, fmt.Errorf("failed to unmarshal native bytes: %v", err)
	}
	return person, nil
}
