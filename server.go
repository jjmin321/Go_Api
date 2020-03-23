package main

import (
	"MaskInfo_api/api"
	"net/http"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	// jwt "github.com/dgrijalva/jwt-go"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `Time = ${time_rfc3339}], Status = ${status}, Method = ${method}, Uri = ${uri}, User_agent = ${user_agent}, Error = ${error},  Remote_ip = ${remote_ip}` + "\n",
	}))
	e.GET("/drugstore", func(c echo.Context) error {
		name := c.QueryParam("name")
		slice1 := []interface{}{}
		ch := make(chan interface{}, 50)
		var wg sync.WaitGroup
		for i := 1; i <= 54; i++ {
			wg.Add(1)
			go api.Drugstore(name, i, &wg, ch)
		}
		wg.Wait()
		close(ch)
		for i := range ch {
			slice1 = append(slice1, i)
		}
		return c.JSON(http.StatusOK, slice1)
	})

	e.GET("/masks", func(c echo.Context) error {
		code := c.QueryParam("code")
		ch := make(chan interface{}, 1)
		var wg sync.WaitGroup
		for i := 1; i <= 51; i++ {
			wg.Add(1)
			go api.Masks(code, i, &wg, ch)
		}
		wg.Wait()
		return c.JSON(http.StatusOK, <-ch)
	})
	e.Logger.Fatal(e.Start(":3000"))
}
