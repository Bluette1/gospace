package main

import "fmt"

func main() {
	greeting := "Hello"
	name := "Stacey"
	sayGreeting(&greeting, &name); // passing arguments by reference

	//Variadic variables
	sum := sum(1, 2, 3, 4, 5)
	fmt.Println(*sum);

		//Methods are like functions that provide a context that a function is
	//executing in
	g := greeter {
		greeting: "Hello",
		name: "Junior",
	}

	g.greet()

	// fmt.Println(divide(1, 0.0));
  d, err := divide(1.0, 0.0) // Comma dilimited list of return values
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d)

	//anonymous functions
	func(){
		fmt.Println("Hello Go")
	}() //invoke the function call

	//In asyncronous code, it is safer to pass the args to the inner scope
  for i := 0; i <= 10; i++ {
		func(i int){
			fmt.Println(i)
		}(i)
	}

	// You can declare a function as a variable. 
	//Functions are types
	var div func(float64, float64) (float64, error)
	div = func(a, b float64) (float64, error) {
		if b == 0.0 {
			return 0.0, fmt.Errorf("Cannot divide by zero");
		}
		return a / b, nil
	}

	fmt.Println(div(6.0, 2.0)) //must be defined before they are executed
}

//parameters of the same type: type is listed once
func sayGreeting(greeting, name * string) { 
	fmt.Println(*greeting, *name);
	*name = "Ted"
	fmt.Println(*name)
}

func sum(values... int) *int { // You can only have one and it has to be at the end
  result := 0
	for _, v := range values {
		result += v
	}
	// fmt.Println(result)
	return &result  // can return a pointer to value
}

	//Handling possible errors
	func divide(a, b float64) (float64, error) {
		if b == 0.0 {
			return 0.0, fmt.Errorf("Cannot divide by zero");
		}
		return a / b, nil
	}
	type greeter struct {
		greeting string
		name string
	}

	//Value receiver
	func (g greeter) greet() { // you can also pass a pointer receiver: func (g *greeter) greet() { 
		fmt.Println("I was called")
		fmt.Println(g.greeting, g.name)
	}