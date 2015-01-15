package main

import (
	"fmt"
	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
	"log"
	"net"
)

const SERVER_ADDR = "127.0.0.1:9323"

type (
	Command struct {
		FuncName   string
		Args       []interface{}
		StatusChan libchan.Sender
	}

	Result struct {
		Result []interface{}
		Status int
	}
)

func (c Command) Sum(a, b int64) int64 {
	fmt.Printf("[sum] %d + %d = %d\n", a, b, a+b)
	return a + b
}

func main() {
	client, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	transport, err := spdy.NewClientTransport(client)
	if err != nil {
		log.Fatal(err)
	}

	sender, err := transport.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}

	receiver, remoteSender := libchan.Pipe()

	command := &Command{
		FuncName:   "Sum",
		Args:       []interface{}{1, 2},
		StatusChan: remoteSender,
	}

	err = sender.Send(command)
	if err != nil {
		log.Fatal(err)
	}

	response := &Result{}
	receiver.Receive(&response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Result: %v", response.Result)
}
