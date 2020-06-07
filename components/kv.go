package components

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"

	"github.com/gomodule/redigo/redis"
)

var mutex sync.Mutex

func ConnectKv() (redis.Conn, error) {
	conn, err := redis.Dial("tcp", "localhost:6379")
	// redis.DialUsername("username"),
	// redis.DialPassword("password"),
	return conn, err
}

// 设置最大值
func SetMax(key string, value int) error {
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return err
	}
	_, ers := conn.Do("SET", key, value)
	return ers
}

// 获取最大值
func GetMax(key string) (int, error) {
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return 0, err
	}
	r, err := conn.Do("GET", key)
	if r == nil || r == 0 {
		return 0, err
	}
	value, err := strconv.Atoi(string(r.(interface{}).([]uint8)))
	if err != nil {
		return 0, errors.New("")
	}
	return value, err
}

// 查看余量
func Surplus(key string) (int, error) {
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return 0, err
	}
	r, err := conn.Do("LLEN", key)
	if err != nil {
		return 0, err
	}
	return int(r.(int64)), nil
}

// 分批写入 写入数不低于公共数字单位，且要求为公共数字单位的整数倍
func PushPipeline(supplement int) error {
	conn, err := ConnectKv()
	defer conn.Close()
	current, _ := GetMax(MaxKey)
	if err != nil {
		return err
	}
	batch := int(math.Floor(float64(supplement) / Unit))
	for i := 0; i < batch; i++ {
		for x := (current + 1); x < (current + int(Unit) + 1); x++ {
			sequence := Generate(int64(x)) // 预生成 预存
			conn.Send("LPUSH", ListKey, int(sequence))
		}
		conn.Flush()
		current = current + int(Unit)
	}
	err = SetMax(MaxKey, current) // 将最大值写入 kv
	return err
}

// 分批读取 读取数不低于公共数字单位 性能约为单次读取的十倍
func RpopPipeline(channel chan int, need int) error {
	mutex.Lock()
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return err
	}
	// 获取远端存储当前余量 用于取值定位
	llen, err := conn.Do("LLEN", ListKey)
	if err != nil {
		return err
	}
	spls := int(llen.(int64))
	fmt.Println("channel len:", len(instance.Channel))
	fmt.Println("kv surplus: ", spls, "  need: ", need)

	// 事务确保批量取值和批量减值正常进行 此操作相当于批量弹出
	conn.Send("MULTI")
	conn.Send("LRANGE", ListKey, spls-need, spls)
	conn.Send("LTRIM", ListKey, spls-need, -1)
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		mutex.Unlock()
		return err
	}

	var sli []int // 便于排序

	// 返回多值 通过控制流和断言进行值的转换
	for _, value := range values {
		switch value.(type) {
		case string:
			continue
		case interface{}:
			for _, val := range value.([]interface{}) {
				if val != nil {
					mst, _ := strconv.Atoi((string(val.([]uint8))))
					sli = append(sli, mst)
				}
			}
		default:
			continue
		}
	}
	// 升序排序
	sort.Slice(sli, func(i, j int) bool {
		return sli[i] < sli[j]
	})
	fmt.Println("补充了 ", len(sli))
	// 按序推入信道
	for _, mst := range sli {
		channel <- mst
	}
	fmt.Println("补充后信道", len(channel))
	mutex.Unlock()
	return err
}

func KvToChannel(channel chan int, need, thresshold int) error {
	err := KvSupplement(thresshold)
	if err != nil {
		return err
	}
	RpopPipeline(channel, need)
	return nil
}

// 补充
func KvSupplement(thresshold int) error {
	mutex.Lock()
	surplus, err := Surplus(ListKey)
	if err != nil {
		mutex.Unlock()
		return err
	}
	if surplus < thresshold {
		fmt.Println("current surplus: ", surplus, "  小于 thresshold: ", thresshold, "  补充一波 ", instance.KvSupplement)
		PushPipeline(instance.KvSupplement)
	}
	mutex.Unlock()
	return nil
}
