package main

import (
	"flag"
	"fmt"
)

func main() {
	testPtr := flag.String("test", "default", "test string")
	flag.Parse()
	fmt.Println(*testPtr)
}
