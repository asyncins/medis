package handler

import (
	"fmt"
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
		fmt.Println("信道容量和当前信道量 ", magazine.Capacity, len(magazine.Channel))
		fmt.Println("inline Freedom:", components.Freedom, "  Need: ", need)
		components.Freedom = 1
		go func() {
			components.KvToChannel(magazine.Channel, need, magazine.KvThreshold)
			components.Freedom = 0
		}()

	}
	number := <-magazine.Channel
	seq := strconv.Itoa(number)
	if len(magazine.Channel)%5000 == 0 {
		fmt.Println("channel len: ", len(magazine.Channel))
	}
	if len(magazine.Channel) < 3 {
		fmt.Println("channel len: ", len(magazine.Channel))
	}
	mutex.Unlock()
	return ctx.String(200, seq)
}
