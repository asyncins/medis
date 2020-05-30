/*
* 薄雾算法
*
* 1      2                                                     48         56       64
* +------+-----------------------------------------------------+----------+----------+
* retain | increas                                             | salt     | sequence |
* +------+-----------------------------------------------------+----------+----------+
* 0      | 0000000000 0000000000 0000000000 0000000000 0000000 | 00000000 | 00000000 |
* +------+-----------------------------------------------------+------------+--------+
*
* 0. 最高位，占 1 位，保持为 0，使得值永远为正数；
* 1. 高位数，占 47 位，高位数（必须是自增数）在高位能保证结果值呈递增态势，遂低位可以为所欲为；
* 2. 随机因子一，占 8 位，上限数值 255，使结果值不可预测；
* 3. 随机因子二，占 8 位，上限数值 255，使结果值不可预测；
*
* 编号上限为百万亿级，上限值计算为 140737488355327 即 int64(1 << 47 - 1)，假设每天取值 10 亿，能使用 385+ 年
 */

package mist

import (
	"fmt"
	"math/rand"
	"medis/common"
	"time"
)

const saltBit = uint(8)                       // 随机因子二进制位数
const sequenceBit = uint(8)                   // 序列号二进制位数
const saltMax = int64(1<<saltBit - 1)         // 随机因子数值上限
const sequenceMax = int64(1<<sequenceBit - 1) // 序列号数值上限
const saltShift = sequenceBit                 // 随机因子移位数
const increasShift = saltBit + sequenceBit    // 高位数移位数

/* 生成区间范围内的随机数 */
func RandInt64(duration int64) int64 {
	rand.Seed(duration) // 时间戳作为随机因子
	return rand.Int63n(common.RandMax)
}

/* 生成唯一编号 */
func Generate(increas int64) int64 {
	now := time.Now()
	saltA := RandInt64(now.UnixNano() / 1e1)
	saltB := RandInt64(now.UnixNano() / 1e2)
	fmt.Println(saltA, saltB)
	mist := int64((increas << increasShift) | (saltA << saltShift) | saltB) // 通过位运算实现自动占位
	return mist
}
