package converters

import (
	"encoding/xml"
	"github.com/golang/protobuf/proto"
)

type StringMap map[string]string

type Person struct {
	XMLName xml.Name `xml:"Person"`

	Name     string    `xml:"name" protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Age      int32     `protobuf:"varint,2,opt,name=age,json=age" json:"age,omitempty" xml:"age"`
	Siblings StringMap `protobuf:"bytes,3,rep,name=siblings,proto3" json:"siblings,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Cars     []string  `xml:"cars" protobuf:"bytes,4,rep,name=cars,proto3" json:"cars,omitempty" `
}

// proto settings
func (m *Person) Reset()         { *m = Person{} }
func (m *Person) String() string { return proto.CompactTextString(m) }
func (m *Person) ProtoMessage()  {}

type Converter interface {
	Serialize(p *Person) ([]byte, error)
	Deserialize(bytes []byte) (*Person, error)
}
