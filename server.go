package main

import (
	"medis/components"
	"medis/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()
	server.GET("/seqence", handler.Seqence)
	current, err := components.GetMax(components.MaxKey)
	if err != nil {
		server.Logger.Fatal(server.Start(":1323"))
	}
	if current == 0 {
		components.MagazineInstance(false)
		components.SetMax(components.MaxKey, components.Capacity)
	} else {
		components.MagazineInstance(true)
	}
	server.Logger.Fatal(server.Start(":1323"))
}
