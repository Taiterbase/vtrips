package utils

import "fmt"

func MakeKey(field, value string) []byte {
	return []byte(fmt.Sprintf("%s:%s", field, value))
}
