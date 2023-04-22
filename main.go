package main

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"serialization_tester/converters"
	"serialization_tester/converters/XML"
	"serialization_tester/converters/avro"
	"serialization_tester/converters/json"
	"serialization_tester/converters/msgpack"
	"serialization_tester/converters/native"
	"serialization_tester/converters/proto"
	"serialization_tester/converters/yaml"
)

type Controller struct {
	port int32

	xmlConverter     *XML.Converter
	nativeConverter  *native.Converter
	jsonConverter    *json.Converter
	protoConverter   *proto.Converter
	avroConverter    *avro.Converter
	yamlConverter    *yaml.Converter
	msgpackConverter *msgpack.Converter
}

func NewController(port int32) (*Controller, error) {
	ctrl := &Controller{
		port:             port,
		xmlConverter:     &XML.Converter{},
		nativeConverter:  &native.Converter{},
		jsonConverter:    &json.Converter{},
		protoConverter:   &proto.Converter{},
		avroConverter:    &avro.Converter{},
		yamlConverter:    &yaml.Converter{},
		msgpackConverter: &msgpack.Converter{},
	}
	err := ctrl.avroConverter.SetSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to set schema: %v", err)
	}
	return ctrl, nil
}

func (c *Controller) ProcessConverter(format string, converter converters.Converter) (string, error) {
	person := &converters.Person{
		Name: "Albert",
		Age:  50,
		Siblings: map[string]string{
			"Ameli": "shaml",
			"Azali": "shaml",
		},
		Cars: []string{
			"abc", "def", "dgx",
		},
	}
	structSize := reflect.TypeOf(person).Size()

	var totalTimeSerialize int64
	var totalTimeDeserialize int64

	{
		// fictive serialize and deserialize to init all things in root of library
		bytes, _ := converter.Serialize(person)
		converter.Deserialize(bytes)
	}

	var attempts int64 = 1000
	for i := 0; i < int(attempts); i++ {
		start := time.Now()
		bytes, err := converter.Serialize(person)
		totalTimeSerialize += time.Since(start).Microseconds()
		if err != nil {
			return "", fmt.Errorf("failed to serialize string: %v", err)
		}

		start = time.Now()
		_, err = converter.Deserialize(bytes)
		totalTimeDeserialize += time.Since(start).Microseconds()
		if err != nil {
			return "", fmt.Errorf("failed to deserialize string: %v", err)
		}
	}
	return fmt.Sprintf(
		"%s - %d - %dmcs - %dmcs\n",
		format, structSize, totalTimeSerialize/attempts, totalTimeDeserialize/attempts), nil
}

func (c *Controller) ProcessRequest(buf []byte) (string, error) {
	var conv converters.Converter
	convertersMap := map[string]converters.Converter{
		"xml":     c.xmlConverter,
		"native":  c.nativeConverter,
		"json":    c.jsonConverter,
		"proto":   c.protoConverter,
		"avro":    c.avroConverter,
		"yaml":    c.yamlConverter,
		"msgpack": c.msgpackConverter,
	}
	format := strings.Trim(string(buf), "\n")

	if format == "all" {
		var result []string
		for format = range convertersMap {
			formatResult, err := c.ProcessRequest([]byte(format))
			if err != nil {
				return "", fmt.Errorf("failed to process format %q: %v", format, err)
			}
			result = append(result, formatResult)
		}
		return strings.Join(result, ""), nil
	}

	conv, ok := convertersMap[format]
	if !ok {
		return "", fmt.Errorf("unknown format: %s", format)
	}
	res, err := c.ProcessConverter(format, conv)
	if err != nil {
		return "", fmt.Errorf("failed to process converter %q: %v", format, err)
	}
	return res, nil
}

func (c *Controller) Listen() error {
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		return fmt.Errorf("failed to make listen packet: %v", err)
	}
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("failed to read data from connection: %v", err)
			continue
		}
		res, err := c.ProcessRequest(buf[:n])
		if err != nil {
			fmt.Printf("failed to process request: %v", err)
			continue
		}
		to, err := conn.WriteTo([]byte(res), addr)
		if err != nil {
			fmt.Printf("failed to write to addr %v: %v", to, err)
		}
	}
}

func main() {
	port := os.Args[1]
	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	controller, err := NewController(int32(intPort))
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		err := controller.Listen()
		if err != nil {
			panic(err)
		}
		wg.Done()
	}()

	go func() {
		if len(os.Args) < 3 {
			fmt.Println(1233122)
		}
		format := os.Args[2]
		result, err := controller.ProcessRequest([]byte(format))
		if err != nil {
			panic(fmt.Errorf("failed to get result from format %q: %v", format, err))
		}
		wg.Done()
		fmt.Println(result)
	}()
	wg.Wait()
}
