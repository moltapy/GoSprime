package main

import (
	"fmt"
	parse "gosprime/PreData/Parse"
)

var args parse.Args

func init() {
	args = parse.Args{}
	args.Parse()
}

func main() {
	fmt.Println(*args.WorkPath)
}
