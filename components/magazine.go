package components

import "sync"

var instance *Magazine
var once sync.Once

type Magazine struct {
	ListKey    string   // 存储已生成数据的键名
	MaxKey     string   // 存储当前最大值的键名
	Capacity   int      // 信道容量
	Threshold  int      // 信道容量阈值
	Persent    float64  // 信道容量阈值百分比
	Supplement int      // 补充量
	Channel    chan int // 信道
	Batch      int      // 每批数量
}

func MagazineInstance(empty bool) *Magazine {
	once.Do(func() {
		capacity := Capacity
		channel := make(chan int, capacity)
		threshold := int(float64(capacity) * Persent)
		// 初始化薄雾结构体时生成第一批值，同时将最大值存入 kv
		if empty == false {
			for i := 1; i < capacity+1; i++ {
				channel <- i
			}
		}
		instance = &Magazine{
			ListKey: ListKey, MaxKey: MaxKey, Capacity: capacity,
			Threshold: threshold, Persent: Persent, Supplement: capacity * Multiple,
			Channel: channel, Batch: capacity,
		}
	})
	return instance
}

func GetMagazine() *Magazine {
	return instance
}
