package kv

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

func newPoolFunc() (redis.Conn, error) {
	return redis.Dial("tcp", ":6379")
}

func CreatePool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		Dial:        newPoolFunc,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
	}
}

// 设置最大值
func SetMax(conn redis.Conn, key string, val int) error {
	_, err := conn.Do("SET", key, val)
	return err
}

// 获取最大值
func GetMax(conn redis.Conn, key string) (int, error) {
	r, err := conn.Do("GET", key)
	if r == nil {
		return 0, errors.New("")
	}
	last, _ := strconv.Atoi(string(r.(interface{}).([]uint8)))
	return last, err
}

// 查看余量
func Surplus(conn redis.Conn, key string) (int, error) {
	r, err := conn.Do("LLEN", key)
	if err != nil {
		return 0, err
	}
	return int(r.(int64)), nil
}

// 移除列表最右元素
func Rpop(conn redis.Conn, key string) (int, error) {
	var value int
	r, err := conn.Do("RPOP", key)
	if err != nil || r == nil {
		return value, errors.New("")
	}
	value, _ = strconv.Atoi(string(r.(interface{}).([]uint8)))
	return value, nil
}

// 向列表最左添加元素
func Lpush(conn redis.Conn, key string, val int) error {
	_, err := conn.Do("LPUSH", key, val)
	return err
}

// 分批写入
func PushPipeline(pol *redis.Pool, supplement, current int, lk, mk string) error {
	// 写入数不低于万，且要求为万的整数倍
	conn := pol.Get()
	defer conn.Close()
	batch := int(math.Floor(float64(supplement) / 1e4))
	for i := 0; i < batch; i++ {
		for x := (current + 1); x < (current + 1e4 + 1); x++ {
			conn.Send("LPUSH", lk, x)
		}
		conn.Flush()
		current = current + 1e4
	}
	err := SetMax(conn, mk, current) // 将最大值写入 kv
	return err
}
