package handler

import (
	"medis/components"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var mutex sync.RWMutex
var lock sync.Mutex

func Seqence(ctx echo.Context) error {
	mutex.Lock()
	var magazine = components.GetMagazine()
	// 取值时从信道取值，每次判断信道余量，余量小于等于阈值时从 kv 拉取一批值
	if len(magazine.Channel) <= magazine.Threshold && components.Freedom == 0 {
		need := magazine.Capacity - len(magazine.Channel)
		// 用全局变量锁定 Gorouting 用完再释放
		components.Freedom = 1
		go func() {
			components.KvToChannel(magazine.Channel, need, magazine.KvThreshold)
			components.Freedom = 0
		}()
	}
	number := <-magazine.Channel
	seq := strconv.Itoa(number)
	mutex.Unlock()
	return ctx.String(200, seq)
}
