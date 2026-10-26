package main

import "fmt"

type MyStruct struct {
	Value string
}

type MyInterface interface {
	GetValue() string
}

func (m *MyStruct) GetValue() string {
	return m.Value
}

func ToNil(obj MyInterface) {
	if obj == nil {
		fmt.Printf("obj = nil")
		return
	}
	obj = nil
}

func main() {
	// test
	// nil
	myStruct := &MyStruct{Value: "test"}

	ToNil(myStruct)
	fmt.Printf("myStruct: %+v \n", myStruct)

	myStruct = nil

	ToNil(myStruct)
	fmt.Printf("myStruct: %+v \n", myStruct)
}
