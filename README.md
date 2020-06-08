# Medis 高性能的全局唯一 ID 发号服务

Medis 是薄雾算法 Mist 的工程实践，其名取自 Mist 和 Redis。[薄雾算法](https://github.com/asyncins/mist/blob/master/README.md)是一款性能强到令我惊喜的全局唯一 ID 算法，我将它与业内同样高性能的 Redis 和 Golang 结合到一起，碰撞出了 TPS 为 2.5w/sec 这样超高性能的工程。

有了 Mist 和 Medis，你就拥有了和[美团 Leaf](https://tech.meituan.com/2017/04/21/mt-leaf.html)、[微信 Seqsvr](https://www.infoq.cn/article/wechat-serial-number-generator-architecture)、[百度 UIDGenerator](https://github.com/baidu/uid-generator) 性能相当（甚至超过）的全局唯一 ID 服务了。相比复杂的 UIDGenerator 双 Buffer 优化和 Leaf-Snowflake，薄雾算法 Mist 简单太多了。

### 分布式环境下的 CAP 选择？

你可以基于 Mist 算法打造一个 CP 或者 AP 的分布式服务架构，因为 Mist 足够简单，你只需要专心设计 CAP 的取舍即可。

无数的架构师说过，**分布式环境下的服务，简单即是强**。


### Medis 的性能测试

选用 Jmeter 作为性能测试工具，测试参数：`1 秒内启动100 并发无限循环`，即

```
Number of Threads  100
Ramp-up period     1
Loop Count         Infinte    
```

在 Jmeter 测试的同时，我还用 Gorouting * 10 开启了对数据正确性的测试，数据正确无误；

小数值（十万级）测试会导致频繁的预读预写 Redis 和 Channel，性能约为 2w/sec TPS

大数值（百万级）测试贴近真是业务环境，以下是 Channel 容量为 500W 时的测试结果

| Samples  | Average  |  Median |    90% Line  |   95% Line  |   99% Line  |   Throughput  |
|  ----    |   ----   |  ----   |      ----    |     ----    |     ----    |      ----     |
|  30000000  |    3    |    2   |       4     |      6     |      12    |     27070/sec  |
|  50000000  |    4    |    3    |       7     |      10    |      21    |     22369/sec  |


### Medis 的高性能从何而来？

作为开发者，你一定想知道 Medis 这 2.5w/sec 的 TPS 到底从何而来。实际上这不仅是薄雾算法本身超高性能带来的结果，我在设计上也做了很多尝试，例如：

- 使用 Channel 作为**数据缓存**，这个操作使得发号服务性能提升了 7 倍；
- 采用**预存预取**的策略保证 Channel 在大多数情况下都有值，从而能够迅速响应客户端发来的请求；
- 用 **Gorouting** 去执行耗费时间的预存预取操作，不会影响对客户端请求的响应；
- 采用 **Lrange Ltrim 组合**从 Redis 中批量取值，这比循环单次读取或者管道批量读取的效率更高；
- 写入 Redis 时采用**管道**批量写入，效率比循环单次写入更高；
- Seqence 值的计算**在预存前进行**，这样就不会耽误对客户端请求的响应，虽然薄雾算法的性能是纳秒级别，但并发高的时候也造成一些性能损耗，放在预存时计算显然更香；
- 得益于 Golang Echo 框架和 Golang 本身的高性能，整套流程下来我很满意，如果要追求极致性能，我推荐大家试试 Rust；


### 预存预取是什么流程？

预存预取是 Medis 高性能的基础之一，待我补充流程图。

### 致谢

谢谢 @青南 在 Redis 读取环节提供的 Lrange Ltrim 建议，我替换掉了循环 RPOP 操作，这使得读性能飙升；

谢谢 @Manjusaka 和 @夏溪辰 在 Gorouting 触发环节提供的全局变量和定时器建议，这里选用了全局变量锁定 Gorouting，效果相当好；

谢谢 @夜幕团队 @崔庆才 @大鱼 @Lock 在性能测试和 Golang 方面的建议；
