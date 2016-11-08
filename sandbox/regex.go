package main

import (
	"fmt"
	"regexp"
)

func main() {
  /*
	re := regexp.MustCompile("a(x*)b")
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "T"))
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1"))
	fmt.Println(re.ReplaceAllString("-ab-axxb-", "${1}"))
  */

  re := regexp.MustCompile(`\*\*(.*)\*\*`)
  fmt.Println(re.ReplaceAllString("this is **bold** text", "XXX${1}XXX"))
  fmt.Println(re.ReplaceAllString("-ab-axxb-", "$1"))
  fmt.Println(re.ReplaceAllString("-ab-axxb-", "${1}"))


}
