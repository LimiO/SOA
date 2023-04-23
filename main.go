package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"

	"serialization_tester/converters"
	"serialization_tester/converters/XML"
	"serialization_tester/converters/avro"
	"serialization_tester/converters/json"
	"serialization_tester/converters/proto"
	"serialization_tester/proxy"
)

type Server interface {
	ProcessRequest([]byte) (string, error)
}

type MutlicastServer interface {
	ListenMulticastGroup() error
}

type Controller struct {
	groupAddr string
	port      int32
	server    Server
}

func NewController(port int32, serverType string) (*Controller, error) {
	groupAddr := os.Getenv("GROUP_ADDR")
	if groupAddr == "" {
		return nil, fmt.Errorf("GROUP_ADDR cannot be empty")
	}

	var server Server
	switch serverType {
	case "proxy":
		server = proxy.Server{
			Port: port,
			ConvertersAddrs: map[string]string{
				"xml": ":3000",
				//"proto": ":3001",
				//"json":  ":3002",
			},
			MulticastAddr: groupAddr,
		}
	case "xml":
		server = converters.NewServer(3000, "xml", groupAddr, &XML.Converter{})
	case "proto":
		server = converters.NewServer(3001, "proto", groupAddr, &proto.Converter{})
	case "json":
		server = converters.NewServer(3002, "json", groupAddr, &json.Converter{})
	case "avro":
		conv := &avro.Converter{}
		err := conv.SetSchema()
		if err != nil {
			return nil, fmt.Errorf("failed to set schema: %v", err)
		}
		server = converters.NewServer(3003, "proxy", groupAddr, conv)
	}

	ctrl := &Controller{
		port:      port,
		server:    server,
		groupAddr: groupAddr,
	}
	return ctrl, nil
}

func (c Controller) Listen() error {
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
		res, err := c.server.ProcessRequest(buf[:n])
		if err != nil {
			fmt.Printf("failed to process request: %v", err)
			_, err = conn.WriteTo([]byte(err.Error()), addr)
			continue
		}
		to, err := conn.WriteTo([]byte(res), addr)
		if err != nil {
			fmt.Printf("failed to write to addr %v: %v", to, err)
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Errorf("not enought args"))
	}
	port := os.Args[1]
	intPort, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	controller, err := NewController(int32(intPort), os.Args[2])
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = controller.Listen()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer wg.Done()
		obj, ok := controller.server.(MutlicastServer)
		if !ok {
			fmt.Println("SKIP")
			return
		}
		err = obj.ListenMulticastGroup()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
