[TOC]

## RPC 介绍

Remote Procedure Call，远程过程调用)是一种计算机通信协议，允许调用不同进程空间的程序。RPC 的客户端和服务器可以在一台机器上，也可以在不同的机器上。程序员使用时，就像调用本地程序一样，无需关注内部的实现细节。

一个典型的 RPC 调用如下：

```go
client.call("ServerName.ServiceName", args, &reply) error
```

客户端发送的请求包括服务名 ServerName，方法名 ServiceName，参数 args 3 个，服务端的响应包括错误 error，返回值 reply 2 个。