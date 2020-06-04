package main

import (
	"bencode/parser"
	"bytes"
	"fmt"
)

func main() {
	// input := []byte("i42e")
	// input := []byte("3:foo")
	// input := []byte("li42e3:fooe")
	input := []byte("d5:helloi-3e4:spam3:fooe")
	fmt.Println(parser.Decode(bytes.NewReader(input)))
}