package main

import (
	"fmt"
	//"time"
)

func main() {
	var s string
	for {
		fmt.Scanln(&s)

		for _, val := range s {
			fmt.Print(val, "-")
		}

		fmt.Println("")
		fmt.Println([]byte(s))
	}
}
