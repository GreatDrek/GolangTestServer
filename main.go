package main

import (
	"fmt"
	"math/rand"
)

func main() {
	s0 := make([]int32, 0, 10)

	for i := 0; i < 10; i++ {
		s0 = append(s0, rand.Int31())
	}

	for _, val := range s0 {
		fmt.Println(val)
	}
}
