package handler

import "github.com/labstack/echo/v4"

func Index(ctx echo.Context) error{

	return ctx.String(200, "nil")
}