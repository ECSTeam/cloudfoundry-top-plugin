package main

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

func main() {

	expression, err := govaluate.NewEvaluableExpression("foo > 0")

	parameters := make(map[string]interface{}, 8)
	parameters["foo"] = "2"

	result, err := expression.Evaluate(parameters)
	// result is now set to "false", the bool value.

	fmt.Printf("err: %v\n", err)
	fmt.Printf("result: %v\n", result)
}
