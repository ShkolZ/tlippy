package main

import (
	"fmt"

	"github.com/ShkolZ/tlippy/internal/app"
)

func main() {
	app := app.NewCLI()
	fmt.Println("Tlippy CLI")
	app.Run()
}
