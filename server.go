package main

import (
	"medis/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()
	server.GET("/", handler.Index)
	server.GET("/seq", handler.Seq)
	server.Logger.Fatal(server.Start(":1323"))
}
