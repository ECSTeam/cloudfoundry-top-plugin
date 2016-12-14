package util

import "fmt"

const MEGABYTE = (1024 * 1024)

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
	ZB
	YB
)

func (b ByteSize) String() string {
	return b.StringWithPrecision(2)
}

func (b ByteSize) StringWithPrecision(precision int) string {

	format := fmt.Sprintf("%%.%vf", precision)

	switch {
	case b >= YB:
		return fmt.Sprintf(format+"YB", b/YB)
	case b >= ZB:
		return fmt.Sprintf(format+"ZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf(format+"EB", b/EB)
	case b >= PB:
		return fmt.Sprintf(format+"PB", b/PB)
	case b >= TB:
		return fmt.Sprintf(format+"TB", b/TB)
	case b >= GB:
		return fmt.Sprintf(format+"GB", b/GB)
	case b >= MB:
		return fmt.Sprintf(format+"MB", b/MB)
	case b >= KB:
		return fmt.Sprintf(format+"KB", b/KB)
	}
	return fmt.Sprintf(format+"B", b)
}

func main() {
	fmt.Println(YB, ByteSize(1e13))
}
