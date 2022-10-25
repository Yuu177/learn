[TOC]

# go 优雅编程

## 返回结构体还是接口

### 返回接口

```go
package service

type IService interface {
	ListPosts()
}

type service struct {}

func NewService() IService {
    return &service{}
}

func (s *service) ListPosts() {}
```

#### 讨论

面向接口部分，`NewService` 返回一个 `Interface`，这是一个很危险的设计。往简单的看，这样导致了返回的 `struct` 无法被 copy。往复杂了看，这个时候一个 `interface` 就变成你的公开接口，这个 `interface` 即使只做了向后兼容的改动（例如新增一个 `method` ）都要 break 很多外部代码。

怎么理解上面的话呢？

就是我们定义了一个接口，这个接口就是公开的，如果别人的结构体实现了该接口来使用。当我们的接口新增一个方法后，那么别人的结构体就会因为没有实现这个接口导致代码就会报错。

### 返回结构体

```go
package service

type Service struct {}

func NewService() Service {
    return &service{}
}

func (s *Service) ListPosts() {}
```

这个时候我们会产生这样子的疑惑：不返回 `interface` 难道返回一个 `struct` 么，那上层依赖的还是 `struct`，为什么还需要面向接口，`NewService` 的作用就是返回接口，隐藏内部的实现结构体。

### 相关资料

我又查了一下相关的资料，发现关于这个问题，很多人写过文章，大部分的文章都在提到了一句话『Accept interfaces return structs』，也就是接受接口并且返回结构体。感兴趣可以阅读这篇文章：[How To Use Go Interfaces](https://blog.chewxy.com/2018/03/18/golang-interfaces/)；而 Dave 曾经也发过一条 `Twitter`，我在这里就直接引用一下：

> \#golang top tip: the consumer should define the interface. If you’re defining an interface and an implementation in the same package, you may be doing it wrong.

消费者应该负责定义接口，如果在一个包中同时定义了接口和实现，那么你**可能**就做错了。

### 思考

这样做在某些情况下确实更好并且合理

1. 如果我们上游只存在一个依赖，那么我们返回公开的 `struct` 就比较有价值，上游可以将返回的结构体方法通过 `interface` 进行隔离，去掉不会使用的方法，但是这就需要我们谨慎地定义当前结构体的公有方法以及变量；
2. 如果上游存在多个依赖，为每一个 `package` 单独创建一个 `interface` 就非常麻烦，我们还是需要在新的 `package` 中创建 `interface` 来封装结构体的方式，但是在这种情况下让下游去返回一个 `interface` 相比之下就更加方便；

在一个常见的项目中，使用 `NewService` 的方式返回一个接口，并没有什么问题，无论是 `struct` 无法被 copy 还是 `interface` 增加了方法会导致 break 外部的代码（还要有其他人实现这个接口）都不会有太大的影响，很多时候只有返回 `interface` 才能真正地让别人使用 `interface`.

### 总结

在生产者中返回结构体并消费者中定义接口是更加合理的，我们也应该这么去做，但是我也不认为在 `NewService` 这种『构造器』中返回接口就一定是有问题的。

- 在第一种情况下，结构体是私有的，`New` 方法会返回接口，当前包会依赖其他 `package` 中定义的接口（当然业务代码中基本都是定义在相同的 `package` 下）；
- 在第二种情况下，结构体时公开的，`New` 方法会返回结构体或者结构体指针；

这两种方式怎么选就看你自己了，我在普通的业务服务中会用第一种方式，在框架和库中会用第二种。

## 参考

- [如何写出优雅的 Go 语言代码](https://draveness.me/golang-101/)
- [How To Use Go Interfaces](https://blog.chewxy.com/2018/03/18/golang-interfaces/)