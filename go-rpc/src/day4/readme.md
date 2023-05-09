# day4
- 添加超时处理机制

- 客户端需要进行超时处理的位置：
  - 与服务端建立连接，导致的超时
  - 发送请求到服务端，写报文导致的超时
  - 等待服务端处理时，等待处理导致的超时（比如服务端已挂死，迟迟不响应）
  - 从服务端接收响应时，读报文导致的超时

- 服务端需要进行超时处理的位置：
  - 读取客户端请求报文时，读报文导致的超时
  - 发送响应报文时，写报文导致的超时
  - 调用映射服务的方法时，处理报文导致的超时

- 对于超时处理，需要利用好time包和channel
  - 通过time.After或Context.WithTimeout定义超时时间
  - 通过channel在异步操作之间传递信号，方便对超时情况进行处理