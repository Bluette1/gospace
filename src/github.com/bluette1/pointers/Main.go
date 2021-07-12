// address-of operator
package main

import "fmt"

func main() {
	var a int = 42
	var b *int = &a //asterik declaring pointer to an integer: prefix type with `*`
	// fmt.Println(a, b)
	a = 27
	fmt.Println(a, *b) //asterik dereferencing: * is infront of point

	*b = 14

	fmt.Println(a, *b) //asterik dereferencing

	c := [3]int{1, 2, 3}
	f:= &c[0]
	d := &c[1]
	e := &c[2]
	// e := &c[2] - 4 // this fails; Go doesn't allow pointer arithmetic
	// although unsafe package is available for this usage
	fmt.Printf("%v %p %p %p\n", c, d, e, f)

	// can be declared without declaring underlying type first
	var ms *myStruct
	fmt.Println(ms) // returns nil
	// ms = &myStruct{foo: 42}
	ms = new(myStruct)
	// (*ms).foo = 42
	//short hand: automatic dereferencing in complex types
	ms.foo = 42
	fmt.Println(ms)
	// fmt.Print((*ms).foo)
	fmt.Println(ms.foo)

	//Zero value for a pointer: &{0}
}

type myStruct struct {
	foo int 
}

//arrays are value types
//slices are reference types
//a slice is a projection of the underlying array, NOT an array itself
//maps are also reference types