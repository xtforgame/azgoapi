package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("XXXXXXX")
	doc := js.Global()
	doc.Call("x", "cool")
}
