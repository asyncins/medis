# Medis 高性能的全局唯一 ID 发号服务


### Medis 的性能测试

Jemeter 100 并发

| Samples  | Average  |  Median |    90% Line  |   95% Line  |   99% Line  |   Throughput  |
|  ----    |   ----   |  ----   |      ----    |     ----    |     ----    |      ----     |
|  600000  |    24    |    17   |       37     |      53     |      270    |     4315/sec  |
|  500000  |    19    |    4    |       69     |      118    |      220    |     5071/sec  |


直接返回递增序列的性能约为 5700/sec，