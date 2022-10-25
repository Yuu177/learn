[TOC]

# go 编程模式

## Functional Options

>功能选项是一种模式，将此模式用于您需要扩展的构造函数和其他公共 API 中的可选参数，尤其是在这些功能上已经具有三个或更多参数的情况下。

参考这篇文章：[Go 设置各种选项的最佳套路](https://segmentfault.com/a/1190000024506839)

下面的内容是根据这篇文章做的总结和补充。

### 闭包实现（推荐）

> 参考源码：https://github.com/asim/go-micro/blob/master/options.go

```go
package server

import (
	"time"

	"go.uber.org/zap"
)

type Server struct {
	addr     string
	port     int
	protocol string
	timeout  time.Duration
	logger   *zap.Logger
}

type Option func(*Server)

func WithProtocol(p string) Option {
	return func(s *Server) {
		s.protocol = p
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func NewServer(addr string, port int, opts ...Option) (*Server, error) {
	srv := Server{
		addr:     addr,
		port:     port,
		protocol: "tcp",
		timeout:  30 * time.Second,
		logger:   zap.NewNop(),
	}

	for _, o := range opts {
		o(&srv)
	}

	return &srv, nil
}
```

### 使用 Option 接口实现

#### 方法一

```go
package server

import (
	"time"

	"go.uber.org/zap"
)

type Server struct {
	addr     string
	port     int
	protocol string
	timeout  time.Duration
	logger   *zap.Logger
}

type options struct {
	protocol string
	timeout  time.Duration
	logger   *zap.Logger
}

type Option interface {
	apply(*options)
}

type protocolOption string

func (p protocolOption) apply(opts *options) {
	opts.protocol = string(p)
}

func WithProtocol(p string) Option {
	return protocolOption(p)
}

type timeoutOption time.Duration

func (t timeoutOption) apply(opts *options) {
	opts.timeout = time.Duration(t)
}

func WithTimeout(t time.Duration) Option {
	return timeoutOption(t)
}

type loggerOption struct {
	Log *zap.Logger
}

func (l loggerOption) apply(opts *options) {
	opts.logger = l.Log
}

func WithLogger(log *zap.Logger) Option {
	return loggerOption{Log: log}
}

func NewServer(addr string, port int, opts ...Option) (*Server, error) {
	options := options{
		protocol: "tcp",
		timeout:  30 * time.Second,
		logger:   zap.NewNop(),
	}

	for _, o := range opts {
		o.apply(&options)
	}

	srv := Server{
		addr:     addr,
		port:     port,
		protocol: options.protocol,
		timeout:  options.timeout,
		logger:   options.logger,
	}

	return &srv, nil
}
```

缺点：代码太冗余

#### 方法二（推荐）

```go
package server

import (
	"time"

	"go.uber.org/zap"
)

type Server struct {
	addr     string
	port     int
	protocol string
	timeout  time.Duration
	logger   *zap.Logger
}

type Option interface {
	apply(*Server)
}

type optFunc func(*Server)

func (f optFunc) apply(s *Server) {
	f(s)
}

func ProtocolOption(p string) Option {
	return optFunc(func(s *Server) { // 强转为 optFunc. 类似 int(a)
		s.protocol = p
	})
}

func TimeoutOption(t time.Duration) Option {
	return optFunc(func(s *Server) {
		s.timeout = t
	})
}

func LoggerOption(log *zap.Logger) Option {
	return optFunc(func(s *Server) {
		s.logger = log
	})
}

func NewServer(addr string, port int, opts ...Option) (*Server, error) {
	srv := Server{
		addr:     addr,
		port:     port,
		protocol: "tcp",
		timeout:  30 * time.Second,
		logger:   zap.NewNop(),
	}

	for _, o := range opts {
		o.apply(&srv)
	}

	return &srv, nil
}
```

## 参考链接

- [GO 编程模式](https://coolshell.cn/articles/series/go%e7%bc%96%e7%a8%8b%e6%a8%a1%e5%bc%8f)

- [GO 编程模式：FUNCTIONAL OPTIONS](https://coolshell.cn/articles/21146.html)
- [Uber Go 语言编码规范](https://github.com/xxjwxc/uber_go_guide_cn#%E5%8A%9F%E8%83%BD%E9%80%89%E9%A1%B9)
- [Go 设置各种选项的最佳套路](https://segmentfault.com/a/1190000024506839)