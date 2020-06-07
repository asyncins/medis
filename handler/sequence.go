package handler

import (
	"fmt"
	"medis/components"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var mutex sync.Mutex
var lck int

func Seqence(ctx echo.Context) error {
	mutex.Lock()
	fmt.Println("outsaid lck ", lck)

	var magazine = components.GetMagazine()
	// 取值时从信道取值，每次判断信道余量，余量小于等于阈值时从 kv 拉取一批值
	if len(magazine.Channel) <= magazine.Threshold && lck == 0 {
		lck = 1
		fmt.Println("inline lck ", lck)
		need := magazine.Capacity - len(magazine.Channel)
		go func() {
			components.KvToChannel(magazine.Channel, need, magazine.KvThreshold)
		}()
		lck = 0
	}
	number := <-magazine.Channel
	seq := strconv.Itoa(number)
	// fmt.Println("channel len: ", len(magazine.Channel))
	mutex.Unlock()
	return ctx.String(200, seq)
}
