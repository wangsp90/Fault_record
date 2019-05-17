package main

import (
	"fmt"
)

type Student struct {
	Name  string
	Age   int
	Score float32
}

func main() {
	var a Student
	if a.Name == "" {
		fmt.Println("yes")
	}
}
