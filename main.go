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

}
