package components

import (
	"errors"
	"log"
	"math"
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

// 移除列表最右元素
func Rpop(key string) (int, error) {
	var value int
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return 0, err
	}
	r, err := conn.Do("RPOP", key)
	if err != nil || r == nil {
		return value, errors.New("")
	}
	value, _ = strconv.Atoi(string(r.(interface{}).([]uint8)))
	return value, nil
}

// 向列表最左添加元素
func Lpush(key string, val int) error {
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return err
	}
	_, err = conn.Do("LPUSH", key, val)
	return err
}

// 分批写入
func PushPipeline(supplement, current int) error {
	// 写入数不低于万，且要求为万的整数倍
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return err
	}
	batch := int(math.Floor(float64(supplement) / 1e4))
	for i := 0; i < batch; i++ {
		for x := (current + 1); x < (current + 1e4 + 1); x++ {
			conn.Send("LPUSH", ListKey, x)
		}
		conn.Flush()
		current = current + 1e4
	}
	err = SetMax(MaxKey, current) // 将最大值写入 kv
	return err
}

// 分批读取 性能约为单次读取的十倍
func RpopPipeline(channel chan int, need int) error {
	// 读取数不低于千
	conn, err := ConnectKv()
	defer conn.Close()
	if err != nil {
		return err
	}
	batch := int(math.Floor(float64(need) / 1e3))
	for i := 0; i < batch; i++ {
		for x := 0; x < 1e3; x++ {
			if len(channel) == cap(channel) {
				break // 信道满时停止写入值
			}
			conn.Send("RPOP", ListKey)
			conn.Flush()
			reply, errs := conn.Receive()
			if errs != nil || reply == 0 || reply == nil {
				break // kv 已取空
			}
			value, _ := strconv.Atoi(string(reply.(interface{}).([]uint8)))
			channel <- value
		}
	}
	return err
}

func KvToChannel(channel chan int, need, supplement int) error {
	mutex.Lock()
	err := KvSupplement(supplement)
	if err != nil {
		mutex.Unlock()
		return err
	}
	RpopPipeline(channel, need)
	mutex.Unlock()
	return nil
}

// 补充
func KvSupplement(supplement int) error {
	surplus, err := Surplus(ListKey)
	if err != nil {
		return err
	}
	if surplus < supplement {
		current, ers := GetMax(MaxKey)
		log.Println("max: ", current)
		if ers != nil {
			return ers
		}
		go func() {
			PushPipeline(supplement, current)
		}()
	}
	return nil
}
