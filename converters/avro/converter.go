package avro

import (
	"fmt"

	"github.com/hamba/avro"

	"serialization_tester/converters"
)

type Converter struct {
	schema avro.Schema
}

func (c *Converter) SetSchema() error {
	//{"name": "Cars", "type": "array", "items": "string"}
	schema, err := avro.Parse(`{
		"type": "record",
		"name": "me",
		"namespace": "org.hamba.avro",
		"fields" : [
			{"name": "Name", "type": "string"},
			{"name": "Age", "type": "int"},
			{"name": "Siblings", "type": {"type":"map", "values": "string"}},
			{"name": "Cars", "type": {"type":"array", "items": "string"}}
		]
	}`)
	if err != nil {
		return fmt.Errorf("failed to parse struct: %v", err)
	}
	c.schema = schema
	return nil
}

func (c *Converter) Serialize(p *converters.Person) ([]byte, error) {
	bytes, err := avro.Marshal(c.schema, p)
	if err != nil {
		return nil, fmt.Errorf("avro marshal failed: %v", err)
	}
	return bytes, nil
}

func (c *Converter) Deserialize(raw []byte) (*converters.Person, error) {
	person := &converters.Person{}
	err := avro.Unmarshal(c.schema, raw, person)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal avro bytes: %v", err)
	}
	return person, nil
}
