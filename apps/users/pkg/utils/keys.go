package utils

import "fmt"

func MakeKey(field, value string) []byte {
	return fmt.Appendf([]byte{}, fmt.Sprintf("%s:%s", field, value))
}
