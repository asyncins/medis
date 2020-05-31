package handler

import (
	"fmt"
	"medis/common"
	"medis/kv"
	"medis/mes"
	"medis/mist"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var capacity int
var lock sync.Mutex

func Seq(ctx echo.Context) error {
	conn := pol.Get()
	defer conn.Close()
	var med = &mes.Medis{}

	// 根据 kv 中是否有最大值来决定初始化薄雾结构体时是生成值还是从 kv 拉取值
	capacity, err := kv.GetMax(conn, common.MaxKey)
	if capacity == 0 || err != nil {
		med = mes.NewMedis()
		kv.SetMax(conn, med.MaxKey, med.Capacity)
	} else {
		med = mes.NewEmptyMedis()
	}

	// 取值时从信道取值，每次判断信道余量，余量小于等于阈值时从 kv 拉取一批值
	if len(med.Channel) <= med.Threshold {
		need := med.Capacity - len(med.Channel)
		go func() {
			mes.GetBatch(med.ListKey, med.MaxKey, med.Channel, need, med.Supplement)
		}()
	}
	number := <-med.Channel
	seq := strconv.Itoa(int(mist.Generate(int64(number))))
	fmt.Println(number, seq)
	return ctx.String(200, seq)
}
