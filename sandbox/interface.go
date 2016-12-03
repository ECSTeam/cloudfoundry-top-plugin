package main

/*
import "fmt"

type displayFunc func(dataRow interface{}) string

type Person struct {
	name string
	age  int
}

func mainX() {
	fmt.Printf("start\n")

	personArray := make([]*Person, 0, 10)
	personArray = append(personArray, &Person{"one", 1})
	personArray = append(personArray, &Person{"two", 2})
	personArray = append(personArray, &Person{"three", 3})

	displayFunc := func(dataRow interface{}) string {
		person := dataRow.(*Person)
		return fmt.Sprintf("name:%v age:%v", person.name, person.age)
	}
	loop(personArray, displayFunc)
}

func loop(data interface{}, displayFunc displayFunc) {

	array := data.([]*Person)

	for i, item := range array {
		value := displayFunc(item)
		fmt.Printf("i: %v  item:%v\n", i, value)
	}

}
*/
