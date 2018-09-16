package main

import (
	"fmt"
)

type allPerson struct {
	index int
	person
}

type person struct {
	name     string
	age      byte
	nickname string
	key      int
}

func (ap allPerson) print() {
	fmt.Println(ap.name)
}

func main() {

	allP := allPerson{1, person{"Drek", 24, "GreatDrek", 2410}}

	allP.print()

	// foo это наше замыкание
	foo := outer()

	// вызов замыкания
	for i := 0; i < 5; i++ {
		fmt.Println(foo(4))
	}

	var x oru = 45
	x.print()
	x.polovina()

}

func outer() func(int) int {
	a := 0

	return func(x int) int {
		a += x
		return a
	}
}

type oru int

func (or oru) print() {
	fmt.Println(or)
}

func (or oru) polovina() {
	fmt.Println(float64(or) / float64(2))
}
