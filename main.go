package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Start programm")

	c0 := make(chan byte)
	c1 := make(chan string)

	go func() {
		defer close(c0)
		for {
			c0 <- 66
			time.Sleep(time.Millisecond * 500)
			break
		}
	}()

	go func() {
		for {
			c1 <- "drek"
			time.Sleep(time.Millisecond * 1250)
			//close(c1)
		}
	}()

	go func() {
		for {
			select {
			case c, open := <-c0:
				if open {
					fmt.Println(c, open)
				}
			case c := <-c1:
				fmt.Println(c)
			default:
				fmt.Println("data")
				time.Sleep(time.Millisecond * 800)
			}
		}
	}()

	fmt.Scanln()
	fmt.Println("End programm")
}
