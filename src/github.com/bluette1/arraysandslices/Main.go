package main

import (
	"fmt"
)

func main() {
	// arr := [3]int{85, 43, 47}
	arr := [...]int{85, 43, 47}

	arrA := [3]string{}

	arrB := arrA //(copying by value, NOT by reference)

	arrC := &arrA //(Use the address operation to copy by reference)
	arrB[0] ="2"
	fmt.Println(arr)
	
	arrA[0] = "r"
	fmt.Println(arrA)
	fmt.Println(arrB)
	fmt.Println(arrC)

	//Fixed size is a limitation so they are used to back slices
	//Almost verything we can do with an array we can do with slices
	s := []int{1, 2, 3}
	l := len(s)
	cp := cap(s)
	fmt.Printf("Length: %v\n", l)
	//capacity vs slice
	fmt.Printf("cap: %v\n", cp)

	// Copied by reference !so take care

	sl := s[:] //copy of all elements
  sla := s[:3] //Much like Python syntax

	fmt.Printf("sl: %v\n", sl)

	fmt.Printf("sla: %v\n", sla)

	//Making slices
	a := make([]int, 3, 100)
	fmt.Printf("Length a: %v\n", len(a))
	fmt.Printf("Capacity a: %v\n", cap(a))

	//Appending to a slice
	a = append(a, 1);
	//can append multiple arguments
	// a = append(a, 2, 3, 4, 5);

	// size keeps doubling in general
	//Concantenation / spread operation like JS
	//The above expressions works the same way as
	a = append(a, []int{2, 3, 4, 5}...)

  // Stack operation
	  // - Pop off the stack, 
	//Shift operation
	// a = a[1:] => [0 0 1 2 3 4 5] 
	c := a[:1]
	fmt.Printf("%v \n", c)
  // b := a[:len(a)] =>= a[:]
	b := a[:len(a) - 1]
fmt.Println(b)
}