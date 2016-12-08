package main

import "fmt"

type displayFunc func(dataRow IData) string

type IData interface {
	Id() string
}

type Person struct {
	name string
	age  int
}

func (p *Person) Id() string {
	return p.name
}

func (p *Person) display() string {
	return "hello"
}

func main() {
	fmt.Printf("start\n")

	//personArray := make([]IData, 0, 10)
	personArray := make([]*Person, 0, 10)
	personArray = append(personArray, &Person{name: "one", age: 1})
	personArray = append(personArray, &Person{name: "two", age: 2})
	personArray = append(personArray, &Person{name: "three", age: 3})

	displayFunc := func(dataRow IData) string {
		//fmt.Printf("%v\n", dataRow)
		person := dataRow.(*Person)
		return fmt.Sprintf("name:%v age:%v", person.name, person.age)
	}

	fmt.Printf("len: %v\n", len(personArray))
	displayDataArray := make([]IData, len(personArray))
	for i, d := range personArray {
		displayDataArray[i] = d
		//fmt.Printf("i=%v, d=%v\n", i, d)
	}

	loop(displayDataArray, displayFunc)
}

func loop(dataArray []IData, displayFunc displayFunc) {

	for i, item := range dataArray {
		value := displayFunc(item)
		fmt.Printf("i: %v  item:%v\n", i, value)
	}

}
