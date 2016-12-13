package util

//import "fmt"
import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Unicode characters: http://graphemica.com/unicode/characters/page/34
	Ellipsis = string('\U00002026')
)

func CaseInsensitiveLess(s1, s2 string) bool {
	// TODO: Find a more efficent way to do this that does not involve obj creation
	return strings.ToUpper(s1) < strings.ToUpper(s2)
}

func Format(n int64) string {
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

func FormatDisplayData(value string, size int) string {
	if len(value) > size {
		value = value[0:size-1] + Ellipsis
	}
	format := fmt.Sprintf("%%-%v.%vv", size, size)
	return fmt.Sprintf(format, value)
}
