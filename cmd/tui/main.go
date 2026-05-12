package main

import (
	"fmt"

	"github.com/ShkolZ/tlippy/internal/app"
)

func main() {
	tui := app.NewTUI()
	_, err := tui.Run()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
