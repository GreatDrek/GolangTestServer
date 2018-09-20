package main

import (
	"fmt"
	//"time"
)

type myInterface interface {
	TestFunc()
}

type myStruct struct {
	name string
}

func (mS myStruct) TestFunc() {
	fmt.Println(mS.name)
}

func main() {
	c := make(chan myInterface)

	go func() {
		mI := myInterface(myStruct{name: "Drek"})
		c <- mI
	}()

	go func() {
		mI1 := <-c
		mI1.TestFunc()
	}()

	fmt.Scanln()

	fmt.Println("End programm")
}
