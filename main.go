package main

import (
	"fmt"
)

type myInterface interface {
	move(int)
	myPrint()
}

type myType struct {
	Number int
}

type myType2 struct {
	Number int
}

func (m *myType) move(x int) {
	m.Number += x
}

func (m *myType2) move(x int) {
	m.Number -= x
}

func (m *myType) myPrint() {
	fmt.Println(m.Number)
}

func (m *myType2) myPrint() {
	fmt.Println(m.Number)
}

func main() {
	fmt.Println("Interface")

	var slice []myInterface = []myInterface{&myType{Number: 4}, &myType2{Number: 4}}
	x := 19
	for _, sl := range slice {
		sl.move(x)
		sl.myPrint()
		x++
	}

	var i I

	i = "d"

	fmt.Println(i)
	_, ok := i.(string)
	fmt.Println(ok)

	infoType(i)
}

type I interface{}

func infoType(i I) {
	_, s := i.(string)
	_, in := i.(int)

	if s {
		fmt.Println("string")
	}
	if in {
		fmt.Println("int")
	}
}
