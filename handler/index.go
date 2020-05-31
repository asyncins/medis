package handler

import (
	"medis/common"
	"medis/kv"

	"github.com/labstack/echo/v4"
)

var pol = kv.CreatePool()

func Index(ctx echo.Context) error {
	conn := pol.Get()
	defer conn.Close()
	max, err := kv.GetMax(conn, common.MaxKey)
	if err != nil {
		return ctx.JSONPretty(400, echo.Map{"max": 0, "surplus": 0}, " ")
	}
	surplus, err := kv.Surplus(conn, common.ListKey)
	if err != nil {
		return ctx.JSONPretty(400, echo.Map{"max": 0, "surplus": 0}, " ")
	}
	return ctx.JSONPretty(400, echo.Map{"max": max, "surplus": surplus}, " ")
}
