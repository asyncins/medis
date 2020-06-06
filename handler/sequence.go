package handler

import (
	"fmt"
	"medis/components"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var mutex sync.Mutex

func Seqence(ctx echo.Context) error {
	mutex.Lock()
	var magazine = components.GetMagazine()
	// 取值时从信道取值，每次判断信道余量，余量小于等于阈值时从 kv 拉取一批值
	if len(magazine.Channel) <= magazine.Threshold {
		need := magazine.Capacity - len(magazine.Channel)
		go func() {
			components.KvToChannel(magazine.Channel, need, magazine.Supplement)
		}()
	}
	number := <-magazine.Channel
	seq := strconv.Itoa(number)
	sequence := components.Generate(int64(number))
	fmt.Println("seq api value: ", number, " -- ", sequence)
	mutex.Unlock()
	return ctx.String(200, seq)
}
