package components

const Unit = 1e5                // 公共数字单位
const ListKey = "medis"         // 存储已生成数据的键名
const MaxKey = "mdx"            // 存储当前最大值的键名
const Capacity = int(50 * Unit) // 信道容量
const Persent = float64(0.7)    // 信道容量阈值比
const Multiple = 5              // 用于计算补充量的倍数
const RandMax = 250             // 随机值右闭值

var Freedom = 0 // 作为 Gorouting 锁
