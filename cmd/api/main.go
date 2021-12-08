package main

import (
	"fmt"
	"red/cmd/api/gin_app"
)

func main() {
	gin_app.Start()
	fmt.Println("Start")
}
