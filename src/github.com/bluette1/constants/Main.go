package main

import "fmt"

func main() {
	const a = 42

	// const a int = 42
	fmt.Printf("%v, %T\n", a, a)

	//Inferring the type
	var b int16 = 4
	fmt.Printf("%v, %T\n", a+b, a+b)
  //Enumerated constants
	//Variable block
	// const (
	// e = iota
	// d = iota
	// c = iota
	// )

	const (
		e = iota
		d
		c
		)
	fmt.Println(e)
	fmt.Println(d)
	fmt.Println(c)

  //Resets iota
	const (
		_ = iota //Right-only variable(we are not going to use it)
	  f
		g
		)
		fmt.Println(d)
		fmt.Println(c)
}