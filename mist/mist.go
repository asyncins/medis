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
	"crypto/rand"
	"fmt"
	"math/big"
)

const saltBit = uint(8)                  // 随机因子二进制位数
const saltShift = uint(8)                // 随机因子移位数
const increasShift = saltBit + saltShift // 高位数移位数

/* 生成唯一编号 */
func Generate(increas int64) int64 {
	randA, _ := rand.Int(rand.Reader, big.NewInt(255))
	saltA := randA.Int64()
	randB, _ := rand.Int(rand.Reader, big.NewInt(255))
	saltB := randB.Int64()
	fmt.Println(saltA, saltB)
	mist := int64((increas << increasShift) | (saltA << saltShift) | saltB) // 通过位运算实现自动占位
	return mist
}
