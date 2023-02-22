# Go介绍
### 学习资料
- 基础\
[Go语言中文API](https://studygolang.com/pkgdoc)\
[各种技术入门](https://learnxinyminutes.com/)
- 项目\
[7days-golang](https://geektutu.com/post/gee.html)\
[dmt-分布式事务管理器](https://dtm.pub/guide/start.html)\
[不错的开源golang项目](https://www.zhihu.com/question/48821269)\
[Go语言值得学吗，有哪些好入门的简单项目？](https://www.zhihu.com/question/437250309)
- 资料\
[Go 语言设计与实现](https://draveness.me/golang/)\
[Go语言高级编程(Advanced Go Programming)](https://chai2010.cn/advanced-go-programming-book/)\
[Go语言面试题](https://www.topgoer.cn/docs/gomianshiti/mianshiti)

### Golang产生的原因
- 合理利用多核多CPU的优势提升软件系统性能
- 需要足够简洁高效的语言减小软件系统的复杂度
- 需要一种兼顾开发速度和运行速度的语言，并且能够解决内存泄露问题

---

### 发展历程
- 2007年开始开发
- 2009年正式发布
- 2015大版本，移除了其中最后的C代码
- 2018年，发布1.10版本
- 2022年，发布1.19版本
- 现在较常用的版本为1.9.2版本

---

### 特点
1. 兼顾静态编译语言的**安全和性能**和动态语言的**开发维护的高效率**，Go=C+Python
2. 从C语言中继承了较多的理念
   - 包括：表达式语法、控制结构、基础数据类型、调用参数传值、指针等
   - Go中的指针相对于C语言有了很多弱化
   - 引入了包的概念，一个Go文件都要属于一个包，而不能单独存在
3. 垃圾回收机制，内存自动回收
4. 天然并发(重要)
   - 从语言层面支持并发，可以简单实现
   - goroutine轻量级线程，可以轻松实现大并发处理，高效利用多核
   - 基于CPS并发模型
5. 引入了管道通信机制，实现goroutine之间的通信
6. 函数可以返回多个值
7. 支持切片slice，延迟defer等