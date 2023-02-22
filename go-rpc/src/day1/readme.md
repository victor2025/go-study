# day1
- 定义了传输结构以及编解码相关函数

- 定义了rpc协议结构
  - rpc协议由三部分组成：Option, Header, Body
  - 组织结构如下：
  ```
    | Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
    | <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
  ```
- 定义了连接创建和处理流程
- 建立了服务提供者demo