package main

import (
	"fmt"
	"lab1/internal/api"
)

func main() {
	fmt.Println("application started")
	api.StartServer()
	fmt.Println("application terminated")
}
