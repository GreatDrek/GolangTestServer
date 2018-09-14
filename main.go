package main

import (
	"fmt"
	"math/rand"
)

func init() {
	fmt.Println("test init")
}

func main() {
	s0 := make([]int32, 0, 10)

	for i := 0; i < 10; i++ {
		s0 = append(s0, rand.Int31())
	}

	for _, val := range s0 {
		fmt.Println(val)
	}

	i := 5
	b := &i
	*b = 6
	fmt.Println(i)
}
