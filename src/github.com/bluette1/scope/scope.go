package main

import (
	"fmt"
	"strconv"
)
// Package level scope
// var i int = 0
var I int = 0 //Exported outside of package 
// because of Pascal case
func main() {
	//Shadowing of the already declared variable
	//Block scope
	var i int = 0
	var j string
	j = strconv.Itoa(i)
	fmt.Printf("%v, %T", j, j)
}