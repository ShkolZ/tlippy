package main

import (
	"fmt"

	"github.com/ShkolZ/tlippy/internal/app"
)

func main() {
	app := app.NewApp()
	if err := app.Run(); err != nil {
		fmt.Println(err)
	}
}
