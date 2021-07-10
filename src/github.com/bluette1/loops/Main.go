package main

import (
	"fmt"
)

func main() {
	// for i, j := 0, 0; i < 5; i = i++ { // this works
		// for i, j := 0, 0; i < 5; i, j = i++, j++ { // this doesn't work
	for i, j := 0, 0; i < 5; i, j = i + 1, j + 2 { //Only 3 loops are allowed
		fmt.Println(i, j)
	}

	//Bad idea to manipulate counters inside loops

	i := 0 //scoped to the function

	for ; i < 5; i++ {
		fmt.Println(i)
	}

	//Go's version of a `while` or `do while` loop
	j := 0
	for ; j < 5; { // the second semi colon if you have the first
		fmt.Println(j)
		j++
	}

	//Syntactic sugar without the semi colons
	k := 0
	for j < 5 {
		fmt.Println(k)
		k++
	}

	//Infinite loops
	m := 0
	for  {
		fmt.Println(m)
		m++
		if m == 5 {
			break
		}

		//Continue keyword
	}

	for i := 0; i < 10; i++ {
		if i % 2 == 0 {
			continue
		}
		fmt.Println(i)
	}

	//Nested/double for loops
	Loop: // Use label to break out of both loops
	for i := 0; i < 10; i++ {
		
		for j := 0; j < 10; j++ {
		
			fmt.Println(i * j)

			//Simply using break statement inside of double loop
			//breaks out of only the first one
			if j == 5 {
				break Loop
			}
		}
	}

	//For range loop
	s := []int{1, 2, 3}
	// for _, v := range s {
	// for k := range s {

	for k, v := range s {
		fmt.Println(k, v) //Prints out key, value pairs
		// fmt.Println(k)
		// fmt.Println(v)
	}
	//The above syntax can be used for arrays, objects, strings and channels
}