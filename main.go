package main

import (
	"fmt"
	"math/rand"
)

func init() {
	fmt.Println("test init")
}

func main() {
	One()
	s0 := make([]int32, 0, 10)

	for i := 0; i < 10; i++ {
		s0 = append(s0, rand.Int31())
	}

	for _, val := range s0 {
		fmt.Println(val)
	}

	p := new(int)

	fmt.Println(p)

	*p = 6
	// p = 6 // Ошибка, так как p имеет тип *int а не int

	fmt.Println(p, *p)
}

func One() {
	fmt.Println("One0")
	defer Two()
	fmt.Println("One1")
}

func Two() {
	fmt.Println("Two0")

	fmt.Println("Two1")
}
