package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//Just like the traditional control structures,
	//defer, panic, recover calls also alter the execution flow of the program
	basicUsage()
	res, err := http.Get("http://google.com/robots.txt")
  defer res.Body.Close() //associate the opening and closing of a resource
	if err != nil {
		log.Fatal(err)
	}
  // Avoid using defer when you are opening and closing many resources
	// in a loop as the closing will not happen until after 
	//the function exits. This leads to memory leaks
	robots, err := ioutil.ReadAll(res.Body) //convert input to bytes

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", robots)
}

func basicUsage() {
	// fmt.Println("first")
	// defer fmt.Println("second")
	// fmt.Println("third")

	// The second statement executes after the main function
	// runs but before it returns

	// defer fmt.Println("first")
	// defer fmt.Println("second")
	// defer fmt.Println("third")

  // Result printed out in LIFO(last in first out order):
	// This is the right way to do it if you think about it
	// third
	// second
	// first

	a := 1

	defer fmt.Println(a) // the value printed out is 1
	a = 2
	fmt.Println("start")
		recoverUsage()
	fmt.Println("end")
}

func panicUsage() {
	// panic used to manage "exceptions"
	// Either we `panic` the application or return an error
	// In `better programming`, possible exceptions are minimised for example
	// opening an inexistent file throws no exception(doesn't panic) but
	// returns an error

	// a := 1
	// b := 0
	// fmt.Println(a / b) // panics
  //  panic("This is an error") Throws an error

	// used in web server

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello world"))
	})
  err := http.ListenAndServe(":8080", nil) //Throws an error if we're already
	//using this port

	if err != nil {
		panic(err.Error())
	}
	//Panics don't have to be fatal(lead to the killing of an application)

}

func recoverUsage() {
	// When a panic occurs, the execution of the immediate function is halted
	// but deferred statements are still executed

	// That is why the recover function is called within the defer.

	fmt.Println("start")

	defer func() { //anonymous function
		if err := recover(); err != nil {
			log.Println("Error: ", err)
			//if you can't recover from an error
			// you can repanic that error
			// panic(err)
		}
	}() //defer takes a function call
	panic("Something bad")
	fmt.Println("end")

	//The `end` statement is not printed because of the panic
	//The logic applies for deeply stacked functions
	// For example if we have another function calling this one:

	// func Main() {
		// fmt.Println("start")
		// recoverUsage()
		// fmt.Println("end")
	// }
	
}