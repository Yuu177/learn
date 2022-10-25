[TOC]

# gRPC

gRPC 是 google 开发的一个 RPC 框架，跨语言，跨平台，基于 Protobuf 序列化协议。简单的说就是：是一个 RPC 框架，使用 Protobuf 序列化数据。

RPC(Remote Procedure Call，远程过程调用)是一种计算机通信协议，允许调用不同进程空间的程序。RPC 的客户端和服务器可以在一台机器上，也可以在不同的机器上。程序员使用时，就像调用本地程序一样，无需关注内部的实现细节。

## Protobuf

### 介绍

protobuf 的全名是 Google Protocol Buffers，是 google 开发的，与语言无关，平台无关，可扩展的系列化结构数据的方法，可用于数据通信，数据存储等。简单点说，就是和 json 或者 xml 类似的**结构数据序列化方法**。其特点：

- 跨平台，跨语言
- 序列化后体积更小，二进制形式，速度更快
- 兼容好，protobuf 的设计有很好的向下或者向上兼容

### 编写

#### .proto 文件

Protobuf 文件以 **.proto **作为文件后缀，在.proto文件里定义好数据结构之后，就可以用工具将这个文件翻译成具体的代码

### 安装

- 安装 protoc

编写完的 `.proto` 文件并不能直接使用，需要用 protoc 工具，将文件翻译成对应的代码。

```bash
brew install protobuf
```

- 安装 protoc-gen-go 插件

安装完之后，就可以进行编译，自动的将写好的 `.proto` 文件转为相对应的代码，由于原生的 protoc 只能生成 c++，python 等代码，要生成 go 语言，需要安装 protoc-gen-go 插件。

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## 踩坑

### map 中使用数组

`repeated` 允许字段重复，对于Go语言来说，它会编译成数组(slice of type)类型的格式。如 `repeated uint32 uids` 对应的 go 语言为 `[]uint32 uids`。

#### 方法一：结构体中定义数组

`map ` 字段不能同时使用 `repeated`。如果要实现像 `repeated` 效果，可以像下面这样子写：

```protobuf
// 错误
message MyRequest {
	map<int64,repeated uint64> values = 1;
}

// 正确
message Uint64Array {
	repeated uint64 uint64s = 1;
}

message MyRequest {
	map<int64,Uint64Array> values = 1;
}
```

注意 Uint64Array 被翻译成代码后是一个结构体

```go
type Uint64Array struct {
	Uint64S []uint64
}
```

#### 方法二：转换成 byte 数组

```protobuf
message MyRequest {
	bytes values = 1; // byte 数组
}
```

业务代码中用 json 转换成 byte 数组后赋值给 MyRequest。grpc server 收到后使用 json 反序列化即可。

```go
m := make(map[int64][]uint64)
jbytes, _ := json.Marshal(m)
```

#### 方案三：grpc 流

> 通过 stream 关键字定义 protobuf 文件

待补充完善

```protobuf
// 普通 RPC
rpc SimplePing(PingRequest) returns (PingReply);

// 客户端流式 RPC
rpc ClientStreamPing(stream PingRequest) returns (PingReply);

// 服务器端流式 RPC
rpc ServerStreamPing(PingRequest) returns (stream PingReply);

// 双向流式 RPC
rpc BothStreamPing(stream PingRequest) returns (stream PingReply);
```

## 参考链接

- [Golang gRPC 入门教程](https://mp.weixin.qq.com/s/ntYd-b0f7YU7wOaWOHGGzQ)

- [Protocol Buffer Basics: Go](https://developers.google.com/protocol-buffers/docs/gotutorial)
- [Protobuf 终极教程](https://colobu.com/2019/10/03/protobuf-ultimate-tutorial-in-go/)

