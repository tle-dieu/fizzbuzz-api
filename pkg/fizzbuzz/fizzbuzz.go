package fizzbuzz

import (
	"bytes"
	"errors"
	"strconv"
)

func FizzbuzzAlgo(d Data) string {
	var response bytes.Buffer
	var write bool

	for i := 1; i < d.Limit; i++ {
		write = false
		if i%d.Int1 == 0 {
			response.WriteString(d.Str1)
			write = true
		}
		if i%d.Int2 == 0 {
			response.WriteString(d.Str2)
			write = true
		}
		if !write {
			response.WriteString(strconv.Itoa(i))
		}
		response.WriteString(" ")
	}
	return response.String()
}

func CheckData(d Data) error {
	if d.Int1 <= 0 {
		return errors.New("int1 should by greater than zero")
	}
	if d.Int2 <= 0 {
		return errors.New("int2 should by greater than zero")
	}
	if d.Str1 == "" {
		return errors.New("str1 is empty")
	}
	if d.Str2 == "" {
		return errors.New("str2 is empty")
	}
	return nil
}
