package main

import (
	"authservice/internal/app"
)

func main() {
	app := app.NewApp()
	app.Run()
}
