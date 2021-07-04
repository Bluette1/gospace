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

	 //numerical types
	 //signed integers
	 //unsigned integers: store only bigger positives(no negatives)
	 var a uint16 = 13
	 fmt.Printf("%v , %T\n", a, a)

	 //Bitwise operatios
	//  - or, and, xor

	 //Bit shifting
	 b := 8 // 2^3
	 fmt.Println(b << 3)
	 fmt.Println(b >> 3)

	 //floating point numbers
	 //float32 float64

	 //Complex numbers exist 
	 c := complex(5, 12) //complex128
	 fmt.Printf("%v , %T\n",c, c)

	 s := "this is a string"
	 fmt.Printf("%v , %T\n",s[2], s[2]) // the string type is an alias for byte, 
	 // ASII/utf-8 encoding
	 //105, []unint8
	 fmt.Printf("%v , %T\n",string(s[2]), s[2])

	 byteS := []byte(s)
	 fmt.Printf("%v , %T\n", byteS, byteS)


	 //type rune int32
	 //Default types are fine most of the time
}