package main

import (
	"fmt"
	//"time"
)

var b bool = false

func main() {
	fmt.Println("Start programm")
	myChan := make(chan string, 3)
	myInputChan := make(chan string)
	go inputConsole(myInputChan)
	go input(myChan, myInputChan)
	go output(myChan)

	for {
		if b == true {
			break
		}
		//fmt.Println(time.Second)
		//time.Sleep(time.Second * 5)
	}

	fmt.Println("End programm")
}

func inputConsole(c chan string) {
	var inStr string
	for {
		fmt.Scanln(&inStr)
		if inStr == "exit" {
			b = true
			break
		}
		c <- inStr
	}
}

func input(c chan string, c1 chan string) {
	for {
		s := <-c1
		c <- s
	}
}

func output(c chan string) {
	for {
		if len(c) == cap(c) {
			for i := 0; i < cap(c); i++ {
				fmt.Println("Chanel:", <-c)
			}
			//fmt.Println(len(c), cap(c))
		}
	}
}
