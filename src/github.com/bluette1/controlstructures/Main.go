package main

import (
	"fmt"
)

func main() {
	pops := make(map[string]int)
	pops = map[string]int{
	 "cityOne": 11199,
	 "cityTwo": 23490,
	 "cityThree": 21490,
	 }

	// Initializer syntax
	 if pop, ok := pops["cityTwo"]; ok {
		 fmt.Println(pop)
	 }
	 //if, switch statements are much like in other languages
	 //Applies short circuiting or lazy evaluation

	 //equality operators when applied to floating point 
	 //computations may not evaluate as expected. So further mathematical
	 //manipulation is required, such as the use of the absolute error:
	//  if a/b - 1 < e {do something}, fine tune e (checking if a == b)

	//Falling through

	//Switch statements 
	//- multiple tests can be evaluated at each level
	//Tag syntax
	switch 4 { //can use initialize here switch i := 2 +3; i{}
	case 1, 2, 3:
		 fmt.Println("first")
		case 4, 5, 6:
			fmt.Println("second")
		default:
			fmt.Println("default")
	}

	//Tagless syntax
	i := 10
	switch { //cases are allowed to overlap, first matching case is evaluated
	case i <= 5:
		 fmt.Println("first")
		case i <= 6:
			fmt.Println("second")
		default:
			fmt.Println("default")
	}
	//Go uses implicit breaks

	j := 1
	switch { //can use the `fallthrough` keyword
	case j <= 5:
		 fmt.Println("first")
		 fallthrough //the next case will execute regardless of whether
		 // it is true or false
		case j <= 6:
			fmt.Println("second")
		default:
			fmt.Println("default")
	}
  // Type case statement
	//the interface type is assignable to any type
	// var k interface{} = 1
	var k interface{} = [2]int{}
	switch k.(type) { //cases are allowed to overlap, first matching case is evaluated
		case int:
		 fmt.Println("first")
		case string:
			fmt.Println("second")
		case [2]int:
			fmt.Println("Array of two integers")
		default:
			fmt.Println("default")
	}

	//Early breaking out of execution by use of break keyword

	var m interface{} = 1
	switch m.(type) { //cases are allowed to overlap, first matching case is evaluated
		case int:
		  fmt.Println("first")
		  fmt.Println("First second line")
		  break //can use logical test to break out
		  fmt.Println("First third line")
		case string:
			fmt.Println("second")
		case [2]int:
			fmt.Println("Array of two integers")
		default:
			fmt.Println("default")
	}
}