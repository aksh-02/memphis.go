package main

import (
	"fmt"
)

func main() {

	mem := newMem()
	err := mem.Connect(&mem, "localhost:6666", "golang", "memphis")
	if err != nil {
		fmt.Println(err.Error())
	}
}
