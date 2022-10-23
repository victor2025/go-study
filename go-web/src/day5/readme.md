# day5
- 中间件的添加和使用
  - 如何决定中间件添加的位置
  - 保障扩展性和易用性的兼顾
  
- 针对分组添加中间件实现统一处理

- 实现方法
  - 每个请求Context可以从自己所处分组中获取所有对应的中间件，并保存在自身中
  - 通过Context的Next方法决定中间件的执行流程并定义操作的顺序
```go
func A(c *Context) {
    part1 // 执行用户handler之前执行
    c.Next() // 继续执行接下来的handler(包括用户handler)
    part2 // 执行用户handler之后执行
}
```