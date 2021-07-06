package main

import (
	"fmt"
)

func main() {
	//Maps are an example of reference types
	//Literal syntax of declaring a make
    // pops := map[string]int{
		// "cityOne": 11199,
		// "cityTwo": 23490,
		// "cityThree": 21490,
		// }
		// fmt.Println(pops)
 //Using make func
 pops := make(map[string]int)
 pops = map[string]int{
	"cityOne": 11199,
	"cityTwo": 23490,
	"cityThree": 21490,
	}
	fmt.Println(pops)

	//Map manipulation

	fmt.Println(pops["cityOne"]) //Much like JS

	//Return order is ot guaranteed
	pops["cityFour"] = 1155366
	fmt.Println(pops)

	delete(pops, "Georgia")

	fmt.Println(pops)
	fmt.Println(pops["Georgia"])

	_, ok := pops["oho"] //check if a key exists
	println((ok))
	_, ok = pops["cityFour"] //check if a key exists
	println((ok))
	println((len(pops)))

	//Structs: gathers information that is related to the variable.
	//They are special because they can contain mixed types
	//Structs are value types
	type Doctor struct {
		number int
		actorName string
		companions []string
	}

	//Field name syntax - safe
	aDoctor := Doctor {
		number: 123,
		actorName:"Lucy",
		companions: []string{"John", "Doe", "Mark"},
	}

	fmt.Println(aDoctor)

	//Drilling down
	fmt.Println(aDoctor.companions)

	//Positional syntax - unsafe
	aDoctorPosnl := Doctor {
		123,
		"Lucy",
		 []string{"John", "Doe", "Mark"},
	}
	fmt.Println(aDoctorPosnl)

	//Anonymous struct for shortlived types such as web responses
	aDoc := struct {name string} {name: "Peter Doc"}
	fmt.Printf("aDoc: %v", aDoc) //Copying is by value
	anotherDoc := aDoc
	anotherDoc.name = "John Doctor"
	fmt.Printf("aDoc: %v\n", aDoc)
	fmt.Printf("anotherDoc: %v\n", anotherDoc)

	//To copy by reference: use the address operator(&)
	anotherDo := &aDoc
	anotherDo.name = "John Doctor"
	fmt.Printf("aDoc: %v\n", aDoc)
	fmt.Printf("anotherDoc: %v\n", anotherDo)

	//Go doesn't support the inheritance model but supports composition

	//A bird has animal characteristics: - by allowing embedding
	// of structs but is completely independent of animal
	type Animal struct {
		Name string
		Origin string
	}
 type Bird struct {
	 Animal
	 Speed float32
	 CanFly bool
 }
	// b := Bird {}
	// b.Name = "Emu"
	// b.CanFly = false
	// b.Origin = "Australia"
	// b.Speed = 987556

	// fmt.Println(b.Name)

	//literal syntax of declaring struct
	b := Bird {
		Animal: Animal{Name: "Emu", Origin: "Australia"},
		Speed: 454.98,
		CanFly: false,
	}
	fmt.Println(b.Name)

	//Tags provide validation syntax:- can be passed to validation frameworks
	//Go fieldnames are typically uppercase, JSON format is lowercase
	// The tags are made visible through a reflection interface and take part in typ
	// e identity for structs but are otherwise ignored.
	 type Member struct {
    Age       int `json:"age,string"`
}
//tells json.Unmarshal() to put the age JSON property,
//  a string, and put it in the Age field of Member, 
// translating to an int. In this example,
//  we pass json two informations, separated with a comma.
}