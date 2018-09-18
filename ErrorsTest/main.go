package main

import (
	"fmt"
	//"time"
)

type myInterface interface {
	MyPrint()
}

type myInt int

func (mp myInt) MyPrint() {
	fmt.Println(mp)
}

func (mp *myInt) Reload() {
	*mp = myInt(0)
}

func main() {
	var x myInt = myInt(5)
	go x.Reload()
	//time.Sleep(100 * time.Millisecond)
	test(x)
}

func test(mI myInterface) {
	mI.MyPrint()
}
