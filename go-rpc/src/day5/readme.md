# day5
- 为rpc框架增加http协议的支持
- http协议主要负责连接的创建，在连接创建后，client和server之间的通信仍然采用rpc协议
- 通过劫持(hijack)http请求以获取对应的连接，将连接交给服务器，以通过rpc协议进行通信
- debug页面可以向用户提供清晰的调用状态