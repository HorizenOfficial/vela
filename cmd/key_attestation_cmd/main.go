package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/communication"
	"github.com/mdlayher/vsock"
)

const (
	msgDelimiter = byte('\n')
)

func run() error {
	serverType := flag.String("servertype", "vsocket", "server type: tcp or vsocket")
	tcpAddr := flag.String("addr", "localhost:12345", "tcp server address")
	vsockCid := flag.Uint("cid", 3, "vsock context id")
	vsockPort := flag.Uint("port", 12345, "vsock port")
	timeout := flag.Uint("timeout", 60, "timeout in seconds for reading response")

	flag.Parse()

	var conn net.Conn
	var err error

	switch *serverType {
	case "tcp":
		fmt.Printf("Connecting to TCP server at %s\n", *tcpAddr)
		conn, err = net.Dial("tcp", *tcpAddr)
	case "vsocket":
		fmt.Printf("Connecting to vsocket server at cid=%d, port=%d\n", *vsockCid, *vsockPort)
		conn, err = vsock.Dial(uint32(*vsockCid), uint32(*vsockPort), nil)
	default:
		return fmt.Errorf("invalid server type: %s. Please use 'tcp' or 'vsocket'", *serverType)
	}

	if err != nil {
		return fmt.Errorf("error connecting to server: %w", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	reqMsg := admin.AdminMessage{
		Type: admin.KeyAttestationRequestMessage,
		Data: nil,
	}

	reqBytes, err := json.Marshal(reqMsg)
	if err != nil {
		return fmt.Errorf("error marshalling request: %w", err)
	}

	reqBytes = append(reqBytes, msgDelimiter)

	_, err = conn.Write(reqBytes)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}

	fmt.Println("Sent KeyAttestationRequestMessage")

	readTimeout := time.Duration(*timeout) * time.Second
	conn.SetReadDeadline(time.Now().Add(readTimeout))

	reader := bufio.NewReader(conn)
	respBytes, err := reader.ReadBytes(msgDelimiter)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var respMsg admin.AdminMessage
	if err := json.Unmarshal(respBytes, &respMsg); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}

	fmt.Println("Received response: ")
	switch respMsg.Type {
	case admin.AdminResponseMessage:
		fmt.Println("Type: AdminResponseMessage")
		responseData, err := json.MarshalIndent(respMsg.Data, "  ", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling response data: %w", err)
		}
		fmt.Printf("Data: %s\n", string(responseData))

	case admin.AdminErrorMessage:
		fmt.Println("Type: AdminErrorMessage")
		var errData communication.ErrorData
		dataBytes, _ := json.Marshal(respMsg.Data)
		if err := json.Unmarshal(dataBytes, &errData); err == nil {
			return errors.New(errData.Message)
		} else {
			return fmt.Errorf("received error response with unparsable data: %v", respMsg.Data)
		}

	default:
		return fmt.Errorf("unknown message type received: %d", respMsg.Type)
	}

	return nil
}


func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
