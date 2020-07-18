package main

import (
	maskinfo "My_api/Maskinfo"
	translate "My_api/Translate"
	"My_api/naturallanguage"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.GET("/drugstore/:name", maskinfo.Drugstore)
	e.GET("/masks/:code", maskinfo.Masks)
	e.GET("/translator", translate.TranslatorPage)
	e.GET("/translate", translate.Translate)
	e.GET("/analyze", naturallanguage.Analyze)
	e.Logger.Fatal(e.Start(":3000"))
}
