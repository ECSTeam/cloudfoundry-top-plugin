package util

//import "fmt"
import "strings"

func CaseInsensitiveLess(s1, s2 string) bool {
  // TODO: Find a more efficent way to do this that does not involve obj creation
  return strings.ToUpper(s1) < strings.ToUpper(s2)
}


/*
func main() {
  fmt.Printf("%v\n", CaseInsensitiveLess("aaa", "bbb"))
}
*/
