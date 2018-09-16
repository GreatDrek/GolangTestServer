package main

import (
	"fmt"
)

type myInt int

type myStruct struct {
	newInt myInt
}

func (i *myInt) printInt() {
	fmt.Println(*i)
}

func (i *myStruct) printInt() {
	i.newInt = 8
	fmt.Println(*i)
}

func main() {
	m := myStruct{5}
	fmt.Println(m)

	fmt.Println(m.newInt)

	m.newInt.printInt()

	m.printInt()

	m.newInt.printInt()
}
