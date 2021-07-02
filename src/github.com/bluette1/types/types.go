package main

import (
	"fmt"
	"strconv"
)

func main() {
	var i int = 0
	var j string
	j = strconv.Itoa(i)
	fmt.Printf("%v , %T\n", j, j)

	//boolean type
	var n bool //zero value is false
	fmt.Printf("%v , %T\n", n, n)
	 n = 1 == 2
	 m := 2 == 2

	 fmt.Printf("%v , %T\n", n, n)
	 fmt.Printf("%v , %T\n", m, m)
}