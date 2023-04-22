package proto

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	"serialization_tester/converters"
)

type Converter struct {
}

func (c *Converter) Serialize(p *converters.Person) ([]byte, error) {
	bytes, err := proto.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("proto marshal failed: %v", err)
	}
	return bytes, nil
}

func (c *Converter) Deserialize(raw []byte) (*converters.Person, error) {
	person := &converters.Person{}
	err := proto.Unmarshal(raw, person)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal proto bytes: %v", err)
	}
	return person, nil
}
