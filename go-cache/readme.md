# Study of go-cache based on geektutu
基于golang的分布式缓存
[7days-golang](https://geektutu.com/post/geecache-day1.html)
[golang-api](https://studygolang.com/pkgdoc)

### 介绍
- 模仿groupcache
- 基于内存的分布式缓存中间件

### 关键点
- 数据的键值对存储
- 内存不足时的淘汰策略
- 并发保护
- 集群扩展