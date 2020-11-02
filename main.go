package main

import (
	"back/app"
	"back/app/config"
)

func main() {
	ginApp := app.CreateGinApp()
	_ = ginApp.Run(config.Config.App.Port)
}
