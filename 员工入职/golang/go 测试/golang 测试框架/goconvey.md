[TOC]

## goconvey 介绍

GoConvey 是一款针对 Golang 的测试框架。

安装

- `go get github.com/smartystreets/goconvey`

## 一个简单的 demo

- calc.go

```go
package main

func Add(a int, b int) int {
	return a + b
}
```

- calc_test.go

```go
package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAdd(t *testing.T) {
	Convey("add 1 + 1 = 2", t, func() {
		So(Add(1, 1), ShouldEqual, 2)
	})
}

```

