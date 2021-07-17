package main

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	var w Writer = ConsoleWriter{}
	w.Write([]byte("Hello Go!"))

	myInt := IntCounter(0)
	var inc Incrementer = &myInt

	for i := 0; i <= 10; i++ {
		fmt.Println(inc.increment())
	}
	var wc WriterCloser = NewBufferedWriterCloser()
	wc.Write([]byte("Hello Youtube listeners. This is a test"))
	wc.Close()

	bwc := wc.(*BufferedWriterCloser) //Converting the interface with pointer
	// bwc = wc.(BufferedWriterCloser) //this fails

	// if you declare interface you have to use a reference
	// if any of the methods implemented use a pointer type, 
	// otherwise you can use a value


	fmt.Println(bwc)
	r, ok := wc.(io.Reader)
	if (ok) {
		fmt.Println(r)
	} else {
		fmt.Println("Conversion failed") //Can't convert buffered writer
		// closer to reader
	}
	//the empty interface
	var myObj interface{} = NewBufferedWriterCloser()
	//Everything can be cast to an interface that has no methods on it,
	// including primitives
	if wc, ok := myObj.(WriterCloser); ok {
		wc.Write([]byte("Hello Youtube listeners. This is a test"))
		wc.Close()

	}

	r, ok = myObj.(io.Reader)
	if (ok) {
		fmt.Println(r)
	} else {
		fmt.Println("Conversion failed") //Can't convert buffered writer
	}

	//Type switch using empty interface
	var i interface{} = 0
	switch i.(type) { //type conversion using the type keyword
	case int: 
		fmt.Println("This is an integer")
	case string: 
		fmt.Println("This is an integer")
	default:
		fmt.Println("I don't know what this is")
	}
}

type Writer interface { 
	Write([]byte) (int, error)
}

type Closer interface {
	Close() error
}

type WriterCloser interface {
	Writer
	Closer
}

type Empty interface {
	//empty interface
}


type BufferedWriterCloser struct {
	buffer *bytes.Buffer
}

type ConsoleWriter struct {

}

func (cw ConsoleWriter) Write(data []byte)(int, error) {
	n, err := fmt.Println(string(data))
	return n, err
}

// You don't have to use structs to implement interfaces,
//you can use any type

type Incrementer interface {
	increment() int
}

type IntCounter int

func (ic *IntCounter) increment() int { //type alias/pointer is passed
	*ic++
	return int(*ic)
}

func (bwc BufferedWriterCloser) Write(data []byte)(int, error) {
	n, err := bwc.buffer.Write(data)
	if err != nil {
		return 0, err
	}

	v := make([]byte, 8)
	for bwc.buffer.Len() > 8 {
		_, err := bwc.buffer.Read(v)

		if err != nil {
			return 0, err
		} 

		_, err = fmt.Println(string(v))

		if err != nil {
			return 0, err
		} 

	}

	return n, nil
}

func (bwc BufferedWriterCloser) Close() error {

	for bwc.buffer.Len() > 0 {
		data := bwc.buffer.Next(8)
		_, err := fmt.Println(string(data))
		if err != nil {
			return err
		} 
	}
	return nil
}

func NewBufferedWriterCloser() *BufferedWriterCloser {
	return &BufferedWriterCloser{
		buffer: bytes.NewBuffer([]byte{}),
	}
}

