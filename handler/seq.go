package handler

import (
	"fmt"
	"medis/kv"
	"medis/mes"
	"medis/mist"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
)

var pol = kv.CreatePool()
var capacity int

var lock sync.Mutex

func Seq(ctx echo.Context) error {
	conn := pol.Get()
	defer conn.Close()
	var med = &mes.Medis{}

	// 初始化薄雾结构体时先生成第一批值，同时将最大值存入 kv
	capacity, err := kv.GetMax(conn, "mdx") // 先判断 kv 是否存储最大值
	if capacity == 0 || err != nil {
		med = mes.NewMedis()
		kv.SetMax(conn, med.MaxKey, med.Capacity)
	} else {
		med = mes.NewEmptyMedis()
	}
	// 取值时从信道取值，每次判断信道余量，余量等于阈值时从 kv 拉取一批值
	if len(med.Channel) <= med.Threshold {

		need := med.Capacity - len(med.Channel)
		go func() {
			mes.GetBatch(pol, med.ListKey, med.MaxKey, med.Channel, need, med.Supplement)
		}()
	}
	number := <-med.Channel
	lock.Lock()
	seq := strconv.Itoa(int(mist.Generate(int64(number))))
	fmt.Println(number, seq)
	lock.Unlock()
	return ctx.String(200, seq)
}
