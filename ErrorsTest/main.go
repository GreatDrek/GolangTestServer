package main

import (
	"fmt"
)

type MyError struct {
	What string
}

// Test
func (e *MyError) Error() string {
	return fmt.Sprintf("%s", e.What)
}

func run() error {
	return &MyError{
		"nil",
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Normal")
	}
}
