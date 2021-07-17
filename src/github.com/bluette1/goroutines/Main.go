package main

import (
	"fmt"
	"runtime"
	"sync"
	// "time"
)

var wg = sync.WaitGroup{}
var counter = 0 //use of wg allows us to declare variable globally
var m = sync.RWMutex{}

func main() { //There's a Go routine executing the main function
	//Go handles concurrency and parallelism more efficientlyy than most other 
	//languages 
	//Go routines provide abstraction over low-level system threading and scheduling
	go hello() // go routine 
	runtime.GOMAXPROCS(100) //tuning variable: too high:
	// scheduler is overworked, greater memory overload 
	fmt.Println(runtime.GOMAXPROCS(-1))// Returns the number of execution threads 
	// that were running previously(by default in this case), 4 os threads by default
	// Each os core processor runs a thread; 4 cores => 4 threads running
	// parallel to one another
	// runtime.GOMAXPROCS(1) //=> Use one core/single processor, perfect concurrency
	// without parallelism

	// Parallelism: when tasks literally run at the same time e.g for a multi
	// core processor

	for i := 0; i <= 10; i++ {
		//In practice however, the frequent locking and unlocking of variables
		// slows down the program. Normal execution without goroutines
		// /threads would be more efficient in this case. 
		//Only use goroutines when concurrency/parallelism is absolutely
		//needed
		wg.Add(2)
		m.RLock() 
		go sayHello()
		m.Lock()
		go increment()
	}
	
	msg :="Hello, Go!"
	//`msg` can be accessed using the concept of closures
	go func() {
		fmt.Println(msg)
	}() //pass in value to avoid a `race condition`
	
  
	go func(msg string) {
		fmt.Println(msg)
	}(msg) //pass in value to avoid a `race condition`

	// time.Sleep(100 *time.Millisecond) //Bad practice in real world applications

	// //Use WaitGroup to synchronize goroutines
	// wg := sync.WaitGroup{}

	wg.Add(1)
	go func(msg string) {
		fmt.Println(msg)
		wg.Done()
	}(msg) //pass in value to avoid a `race condition`
	msg = "Hello, Ruby"
	//Use the race flag to detect a race condition
	//go run/build -race srcfile
 wg.Wait() //Only works properly because we are
 // only synchronizing one goroutine

 //In case of multiple goroutines, synchronization doesn't work so well
 //the Go routines race against each other
 //only one routine is allowed write access at a time
 //Use a mutex
 //As many things as we want to can read at the same time, but only 
 //one can write at a time
 //Nothing can read or write until the writing is done

}

func hello() {
	fmt.Println("Hello, Go!")
}

func sayHello() {
	// m.RLock() //this is still not quite right, we must lock in context
	fmt.Printf("Hello #%v\n", counter)
	m.RUnlock()
	wg.Done()
}

func increment() {
	// m.Lock() //this is still not quite right, we must lock in context
	counter++
	m.Unlock()
	wg.Done()
}

