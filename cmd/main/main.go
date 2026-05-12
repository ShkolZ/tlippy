package main

import (
	"fmt"

	"github.com/ShkolZ/tlippy/internal/app"
)

func main() {
	a := app.NewApp()
	input, err := a.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if input == nil {
		// User cancelled (quit before finishing)
		return
	}

	// TODO: wire up download logic using input
	fmt.Printf("Got input: %+v\n", input)
}

