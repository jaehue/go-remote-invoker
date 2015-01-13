package main

import (
	"fmt"
	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
	"log"
	"net"
	"reflect"
)

type Command struct {
	FuncName   string
	Args       []interface{}
	StatusChan libchan.Sender
}

type Result struct {
	Result []interface{}
	Status int
}

const (
	SERVER_ADDR = "127.0.0.1:9323"
)

func (c Command) Sum(a, b int64) int64 {
	fmt.Printf("[sum] %d + %d = %d\n", a, b, a+b)
	return a + b
}

func main() {
	listener, err := net.Listen("tcp", SERVER_ADDR)
	if err != nil {
		log.Fatal(err)
	}

	tl, err := spdy.NewTransportListener(listener, spdy.NoAuthenticator)
	if err != nil {
		log.Fatal(err)
	}

	for {
		t, err := tl.AcceptTransport()
		if err != nil {
			log.Print(err)
			break
		}

		go func() {
			for {
				receiver, err := t.WaitReceiveChannel()
				if err != nil {
					log.Print(err)
					break
				}

				go func() {
					for {
						command := &Command{}
						err := receiver.Receive(command)
						if err != nil {
							log.Print(err)
							break
						}

						result := command.invoke()

						returnResult := &Result{Result: result, Status: 1}
						err = command.StatusChan.Send(returnResult)
						if err != nil {
							log.Print(err)
						}
					}
				}()
			}
		}()
	}
}

func (c Command) invoke() []interface{} {
	f := reflect.ValueOf(c).MethodByName(c.FuncName)

	var in []reflect.Value
	for _, v := range c.Args {
		in = append(in, reflect.ValueOf(v))
	}

	fmt.Println("Call Command Function")
	results := f.Call(in)

	fmt.Println("Complete Command Function.")
	var returns []interface{}
	for _, v := range results {
		returns = append(returns, v.Interface())
	}
	return returns
}
