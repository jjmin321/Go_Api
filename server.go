package main

import (
	maskinfo "My_api/Maskinfo"
	translate "My_api/Translate"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/drugstore/:name", maskinfo.Drugstore)
	e.GET("/masks/:code", maskinfo.Masks)
	e.GET("/translator", translate.TranslatorPage)
	e.GET("/translate", translate.Translate)
	e.Logger.Fatal(e.Start(":3000"))
}
