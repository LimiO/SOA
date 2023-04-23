package proxy

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Server struct {
	Port            int32
	MulticastAddr   string
	ConvertersAddrs map[string]string
	Result          chan string
}

func (s Server) ListenMulticastGroup() error {
	maddr, err := net.ResolveUDPAddr("udp", s.MulticastAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve udp addr: %v", err)
	}
	conn, err := net.ListenMulticastUDP("udp", nil, maddr)
	if err != nil {
		return fmt.Errorf("failed to listen multicast udp: %v", err)
	}
	defer conn.Close()

	for {
		for range s.ConvertersAddrs {
			buf := make([]byte, 1000)
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				return fmt.Errorf("failed to read from UDP conn: %v", err)
			}
			if string(buf[:n]) == "get_info" {
				continue
			}
			s.Result <- string(buf[:n])
		}
	}
	
	return nil
}

func (s Server) ProcessMulticast() (string, error) {
	maddr, err := net.ResolveUDPAddr("udp", s.MulticastAddr)
	if err != nil {
		return "", fmt.Errorf("failed to resolve udp addr: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, maddr)

	if err != nil {
		return "", fmt.Errorf("failed to set listen conn: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("get_info"))
	if err != nil {
		return "", fmt.Errorf("failed to write to addr: %v", err)
	}
	
	var result []string
	for range s.ConvertersAddrs {
		str := <- s.Result
		result = append(result, str)
	}

	return strings.Join(result, ""), nil
}

func (s Server) ProcessRequest(request []byte) (string, error) {
	req := strings.Trim(string(request), "\n")
	var addr string
	if req == "all" {
		return s.ProcessMulticast()
	}
	addr, ok := s.ConvertersAddrs[req]
	if !ok {
		return "", fmt.Errorf("failed to get converter addr: %s", req)
	}
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return "", fmt.Errorf("failed to set conn with addr %q: %v", addr, err)
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "get_info")
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", fmt.Errorf("failed to write data to addr %q: %v", addr, err)
	}
	buf := make([]byte, 1024)
	n, err := bufio.NewReader(conn).Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read data from connection: %v", err)
	}
	return string(buf[:n]), nil
}
