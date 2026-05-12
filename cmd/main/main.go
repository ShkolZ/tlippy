package main

import (
	"fmt"

	"github.com/ShkolZ/tlippy/internal/app"
)

func main() {
	a := app.NewApp()
	_, err := a.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
