package main

import (
	"fmt"
	"reflect"
)

type Command struct {
	Func   interface{}
	Args   []interface{}
	Result []interface{}
	Done   chan bool
}

func sum(a, b int) int {
	fmt.Printf("[sum] %d + %d = %d\n", a, b, a+b)
	return a + b
}

func main() {
	c := Command{
		Func: sum,
		Args: []interface{}{1, 2},
		Done: make(chan bool),
	}

	c.invoke()

	<-c.Done
	fmt.Printf("Result: %v", c.Result)
}

func (c *Command) invoke() {
	go func() {
		f := reflect.ValueOf(c.Func)

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
		c.Result = returns
		c.Done <- true
	}()
}
