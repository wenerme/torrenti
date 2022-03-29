package testx

import (
	"encoding/json"
	"fmt"
)

func PrintJson(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	NoErr(err)
	fmt.Println(string(b))
}

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
