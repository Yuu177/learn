[TOC]

# Gin

## 介绍

Gin 是一个用 Go (Golang) 编写的 HTTP web 框架。

网上已经很多教学内容了，这里主要讲一下 gin 开发的时候踩到的坑。

## 参数绑定

**定义被绑定的结构体**

```go
type StructName struct {
  Xxx  type  `form:"paramName" binding:"required"`
}
```

**标签说明:**

- `form:"paramName"`: `paramName` 为参数的名称
- `binding:"required"`: 代表字段必须绑定

### 绑定 URL 参数

通过使用函数 `BindQuery` 和 `ShouldBindQuery`，用来**只绑定** `GET` 请求中的 `uri` 参数，如：`/funcName?a=x&b=x 中的 a 和 b`。

### 绑定 JSON

使用函数 `BindJSON` 和 `ShouldBindJSON` 来绑定提交的 `JSON` 参数信息（一般是 POST 请求的 body）。

### binding:"required" 无法接收零值

在 gin 框架中，需要接收前端参数时，参数必填，我们一般添加 `binding:"required"` 标签，这样前端如果不传该参数，gin 框架会自动校验，并报 error。

 gin 的参数校验是基于 `validator` 的，如果给了 `required` 标签，则不能传入零值，比如字符串的不能传入空串，int 类型的不能传入 0，bool 类型的不能传入 false。

有时候我们需要参数必填，而且需要可以传入零值。比如性别 sex 有 0 和 1 表示，0 表示女，1 表示男，而且需要必填。这个时候，我们可以通过定义 int 类型的指针解决该问题。同理，其他类型也是定义指针即可。

总结：gin 框架 `binding:"required"` 无法接收零值，只要把类型定义为指针类型即可接受零值

**相关链接**

- [BindJSON validation failed for a required integer field that has zero value](https://github.com/gin-gonic/gin/issues/737)
- [解决 go gin 框架 binding:"required" 无法接收零值的问题](https://www.cnblogs.com/rainbow-tan/p/15457818.html)

## Gin test

[怎样编写 Gin 的测试用例](https://learnku.com/docs/gin-gonic/1.7/testing/11359)

## 参考

- [Gin 框架中文文档](https://learnku.com/docs/gin-gonic/1.7)

- [Gin框架(六):参数绑定](http://liuqh.icu/2021/05/10/go/gin/6-param-bind/)