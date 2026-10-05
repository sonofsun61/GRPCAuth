package main

import (
	"authservice/internal/app"
)

func main() {
	app := app.InitializeApp()
	app.Run()
}
