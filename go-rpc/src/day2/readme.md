# day2
- 创建异步并发的高性能客户端
- 支持rpc调用的函数需要满足的条件：
  - the method’s type is exported.
  - the method is exported.
  - the method has two arguments, both exported (or builtin) types.
  - the method’s second argument is a pointer.
  - the method has return type error.