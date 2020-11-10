package main

import (
	"back/app"
)

func main() {
	ginApp := app.CreateGinApp()
	_ = ginApp.Run()
}
