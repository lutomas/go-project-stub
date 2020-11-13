package main

import (
	"fmt"

	"github.com/lutomas/go-project-stub/types"
)

func main() {
	version := types.NewVersion("main-app-cli")
	fmt.Printf("Version: %+v\n", *version)
	fmt.Println("hello world CLI")
}
