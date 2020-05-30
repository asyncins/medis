package mes

import (
	"fmt"
	"medis/common"
	"medis/kv"
	"sync"

	"github.com/garyburd/redigo/redis"
)

var lock sync.Mutex

type Medis struct {
	ListKey    string   // 存储已生成数据的键名
	MaxKey     string   // 存储当前最大值的键名
	Capacity   int      // 信道容量
	Threshold  int      // 信道容量阈值
	Persent    float64  // 信道容量阈值百分比
	Supplement int      // kv补充量
	Channel    chan int // 信道
	Batch      int      // 每批数量
}

var instance *Medis
var once sync.Once

func NewMedis() *Medis {
	once.Do(func() {
		capacity := common.Capacity
		channel := make(chan int, capacity)
		threshold := int(float64(capacity) * common.Persent)
		for i := 1; i < capacity+1; i++ {
			channel <- i
		}
		instance = &Medis{
			ListKey: "medis", MaxKey: "mdx", Capacity: capacity,
			Threshold: threshold, Persent: common.Persent, Supplement: capacity * common.Multiple,
			Channel: channel, Batch: capacity,
		}
	})
	return instance
}

func NewEmptyMedis() *Medis {
	once.Do(func() {
		capacity := common.Capacity
		channel := make(chan int, capacity)
		threshold := int(float64(capacity) * common.Persent)
		instance = &Medis{
			ListKey: "medis", MaxKey: "mdx", Capacity: capacity,
			Threshold: threshold, Persent: common.Persent, Supplement: capacity * common.Multiple,
			Channel: channel, Batch: capacity,
		}
	})
	return instance
}

/* 取一批值 */
func GetBatch(pol *redis.Pool, lk, mk string, channel chan int, need, supplement int) error {
	// 从 kv 拉取值的同时检查 kv 余量
	// 如果 kv 余量也低于阈值则从 kv 取出上次存储的最大值 接着生成一批值存入 kv
	// 同时将新的最大值存入 kv
	lock.Lock()
	conn := pol.Get()
	defer conn.Close()
	surplus, err := kv.Surplus(conn, lk)
	if err != nil {
		return err
	}

	// 当 kv 余量小于信道阈值的 5 倍时立即补充
	if surplus < supplement {
		fmt.Println(fmt.Sprintf("kv 余量 %v 小于 kv 余量阈值时 %v，补充 %v 到kv", surplus, supplement, supplement))
		current, _ := kv.GetMax(conn, mk)
		err := kv.PushPipeline(pol, supplement, current, lk, mk)
		if err != nil {
			return err
		}
	}

	for i := 0; i < need; i++ {
		if len(channel) == cap(channel) {
			break // 信道满时停止写入值
		} else {
			r, err := kv.Rpop(conn, lk)
			if err != nil && r == 0 {
				continue // kv 已取空
			}
			channel <- r
		}
	}
	lock.Unlock()
	return nil
}
