[TOC]

## gomock 介绍

当待测试的函数/对象的依赖关系很复杂，并且有些依赖不能直接创建，例如数据库连接、文件 I/O 等。这种场景就非常适合使用 mock/stub 测试。简单来说，就是用 mock 对象模拟依赖项的行为。

gomock 主要包含两个部分：gomock 库和辅助测试代码生成工具 mockgen。

使用如下命令即可安装：

- `go get -u github.com/golang/mock/gomock`
- `go get -u github.com/golang/mock/mockgen`

务必开启 Modules 管理（go.mod），不然使用 mockgen 可能会出现这样的提示：

Loading input failed: source directory is outside GOPATH

## 一个简单的 Demo

### 第一步：写业务代码

```go
// db.go
package main

type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
```

假设 DB 是代码中负责与数据库交互的部分，测试用例中不能创建真实的数据库连接。这个时候，如果我们需要测试 `GetFromDB` 这个函数内部的逻辑，就需要 mock 接口 DB。

### 第二步：使用 mockgen 生成 mock.go 代码

一般传递三个参数。包含需要被 mock 的接口得到源文件 source，生成的目标文件 destination，包名 package。

- `mockgen -source=db.go -destination=db_mock.go -package=main`

#### mockgen 模式补充

`mockgen` has two modes of operation: source and reflect.

##### Source mode

Source mode generates mock interfaces from a source file. It is enabled by using the -source flag. Other flags that may be useful in this mode are -imports and -aux_files.

Example:

```bash
mockgen -source=foo.go [other options]
```

##### Reflect mode

Reflect mode generates mock interfaces by building a program that uses reflection to understand interfaces. It is enabled by passing two non-flag arguments: an import path, and a comma-separated list of symbols.

You can use "." to refer to the current path's package.

Example: 

```bash
mockgen database/sql/driver Conn,Driver

# Convenient for `go:generate`.
mockgen . Conn,Driver
```

**补充：database/sql/driver 就是 import 导包的路径，Conn, Driver 为该包下的接口。**

### 第三步：新建 test.go，写测试用例

```go
package main

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockDB(ctrl)
	// 打桩 DB 的 Get 方法
	m.EXPECT().Get(gomock.Eq("Tom")).Return(100, errors.New("not exist")) //当入参为 tom，返回值为 100, error("not exist")

	// 测试 GetFromDB 函数
	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}

```

- 这个测试用例有 2 个目的，一是使用 `ctrl.Finish()` 断言 `DB.Get()` 被是否被调用，如果没有被调用，后续的 mock 就失去了意义（简单来说，打桩了哪个函数，就必须要调用到）。
- 二是测试方法 `GetFromDB()` 的逻辑是否正确（如果 `DB.Get()` 返回 error，那么 `GetFromDB()` 返回 -1）。
- `NewMockDB()` 的定义在 `db_mock.go` 中，由 mockgen 自动生成。

最终代码结构

```bash
.
├── db.go
├── db_mock.go
└── db_test.go
```

运行测试用例

```bash
$ go test -v
=== RUN   TestGetFromDB
--- PASS: TestGetFromDB (0.00s)
PASS
ok      tpy/mockProject 0.566s
```

**小结**

- 注意：只能对接口进行 mock。
- 使用 mock 之前，先安装 gmock 库和 mockgen 工具。
- 使用 mockgen 对接口生成 mock 代码。
- 打桩被该接口的被调用的方法。
- 测试函数。

## 打桩(stubs)

在上面的例子中，<u>当 `Get()` 的参数为 Tom，则返回 error，这称之为打桩(stub)</u>，有明确的参数和返回值是最简单打桩方式。除此之外，检测调用次数、调用顺序，动态设置返回值等方式也经常使用。

### 参数(Eq, Any, Not, Nil)

```go
m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
m.EXPECT().Get(gomock.Any()).Return(630, nil)
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil) 
m.EXPECT().Get(gomock.Nil()).Return(0, errors.New("nil"))  
```

- `Eq(value)` 表示与 value 等价的值。
- `Any()` 可以用来表示任意的入参。
- `Not(value)` 用来表示非 value 以外的值。
- `Nil()` 表示 None 值

### 返回值(Return, Do, DoAndReturn)

```go
// Return
m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
// Do
m.EXPECT().Get(gomock.Any()).Do(func(key string) {
    t.Log(key)
})
// DoAndReturn
m.EXPECT().Get(gomock.Any()).DoAndReturn(func(key string) (int, error) {
    if key == "Sam" {
        return 630, nil
    }
    return 0, errors.New("not exist")
})
```

- `Return` 返回确定的值。
- `Do` 方法被调用时，要执行的操作。忽略返回值。
- `DoAndReturn` 可以动态地控制返回值。

**补充**

- 实现 gmock 第一次调用返回一个值，第二次调用返回另外一个值。

```go
// 单独写两行
m.EXPECT().Get(gomock.Any()).Return(1, nil)
m.EXPECT().Get(gomock.Any()).Return(2, nil)
```

### 调用次数(Times)

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil).Times(2)
	GetFromDB(m, "ABC")
	GetFromDB(m, "DEF")
}
```

- `Times()` 断言 Mock 方法被调用的次数（默认是一次）。
- `MaxTimes()` 最大次数。
- `MinTimes()` 最小次数。
- `AnyTimes()` 任意次数（包括 0 次）。

### 调用顺序(InOrder)

```go
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() // 断言 DB.Get() 方法是否被调用

	m := NewMockDB(ctrl)
	o1 := m.EXPECT().Get(gomock.Eq("Tom")).Return(0, errors.New("not exist"))
	o2 := m.EXPECT().Get(gomock.Eq("Sam")).Return(630, nil)
	gomock.InOrder(o1, o2)
	GetFromDB(m, "Tom")
	GetFromDB(m, "Sam")
}
```

## 编写可 mock 的代码

写可测试的代码与写好测试用例是同等重要的，如何写可 mock 的代码呢？

- mock 作用的是接口，因此将依赖抽象为接口，而不是直接依赖具体的类。
- 不直接依赖的实例，而是使用依赖注入降低耦合性。

**依赖注入**

> 在软件工程中，依赖注入的意思为，给予调用方它所需要的事物。 “依赖”是指可被方法调用的事物。依赖注入形式下，调用方不再直接指使用“依赖”，取而代之是“注入” 。“注入”是指将“依赖”传递给调用方的过程。在“注入”之后，调用方才会调用该“依赖”。传递依赖给调用方，而不是让让调用方直接获得依赖，这个是该设计的根本需求。

**举个例子**

> 如果一个类A 的功能实现需要借助于类 B，那么就称类 B 是类 A 的依赖，如果在类 A 的内部去实例化类 B，那么两者之间会出现较高的耦合，一旦类 B 出现了问题，类 A 也需要进行改造，如果这样的情况较多，每个类之间都有很多依赖，那么就会出现牵一发而动全身的情况，程序会极难维护，并且很容易出现问题。要解决这个问题，就要把 A 类对 B 类的控制权抽离出来，交给一个第三方去做，把控制权反转给第三方，就称作控制反转（IOC Inversion Of Control）。控制反转是一种思想，是能够解决问题的一种可能的结果，而依赖注入（Dependency Injection）就是其最典型的实现方法。由第三方（我们称作 IOC 容器）来控制依赖，把他通过构造函数、属性或者工厂模式等方法，注入到类 A 内，这样就极大程度的对类 A 和类 B 进行了解耦。

如果 `GetFromDB()` 方法长这个样子

```go
func GetFromDB(key string) int {
	db := NewDB()
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
```

对 DB 接口的 mock 并不能作用于 `GetFromDB()` 内部，这样写是没办法进行测试的。那如果将接口 db DB 通过参数传递到 `GetFromDB()`，那么就可以轻而易举地传入 Mock 对象了。

## mock 错误踩坑

### zsh: command not found: mockgen

在 .zshrc 添加环境变量

```bash
# ~/.zshrc
# GO path
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
export GOPATH="/Users/xxxx/go"
export GOBINPATH="$GOPATH/bin"
```

### unknown embedded interface

mock 嵌套的接口（**接口定义在相同包名，但是不同文件**）的时候会报错：unknown embedded interface xxx

需要用到 `-aux_files` 这个 Flag。

- `-aux_files`: A list of additional files that should be consulted to resolve e.g. embedded interfaces defined in a different file. This is specified as a comma-separated list of elements of the form `foo=bar/baz.go`, where `bar/baz.go` is the source file and `foo` is the package name of that file used by the -source file.

**例子**

> testDemo/my_test 为包路径

```bash
mockgen -source=test01.go -destination=test01_mock.go -package=test -aux_files testDemo/my_test=test02.go
```

- file layout

```bash
.
├── go.mod
├── go.sum
└── my_test
    ├── test01.go
    ├── test01_mock.go
    └── test02.go
```

- go.mod

```go
module testDemo

go 1.18

require github.com/golang/mock v1.6.0
```

- test01.go

```go
package test

type Service interface {
	Get()
	ServiceExp
}
```

test01.go 文件 `Service` 接口中的 `ServiceExp` 是一个 embedded interface。它们是在同一个包名为 test 下。

所以要 mock test01.go 文件 `Service` 接口时候，需要使用 `-aux_files` 标志来声明参考  test02.go 中的 `ServiceExp` 接口。如果有多个 aux_file，用 `,` 隔开。比如

```bash
-aux_files testDemo/my_test=test02.go,test03.go
```

- test02.go

```go
package test

type ServiceExp interface {
	Set()
}
```

相关 Issue: [Mock Embedded Interface in the SAME Package](https://github.com/golang/mock/issues/645)

### aborting test due to missing call(s)

因为打桩的函数一次都没有用到，所以会报错。~~打桩函数如果没有调用次数限制，只需加上 AnyTimes() 即可。~~不建议用 `AnyTime()`，建议修改测试用例打桩代码，需要用到哪些打桩函数再 mock。

### 使用 reflect mode 来 mockgen 第三方包代码

因为引第三方包无法使用 source 指定源文件，所以需要用到反射模式去生成 mock 代码。

文件路径：github.com/golang/mock/sample/user.go

```go
package user

type Index interface {
	Get(key string) interface{}
	GetTwo(key1, key2 string) (v1, v2 interface{})
}

type Embed interface {
	RegularMethod()
	Embedded
	imp1.ForeignEmbedded
}

type Embedded interface {
	EmbeddedMethod()
}

```

使用 mockgen 生成 Index, Embed, Embedded 这个三个接口的 mock 代码。

```bash
mockgen -destination mock_user_test.go -package user_test github.com/golang/mock/sample Index,Embed,Embedded
```

如果有 `vendor` folder，那么需要加上 `-build_flags=--mod=mod`

```bash
mockgen -destination mock_user_test.go -package user_test -build_flags=--mod=mod github.com/golang/mock/sample Index,Embed,Embedded
```

相关 Issue: [Unable to generate mocks in reflect mode](https://github.com/golang/mock/issues/494)

### EXPECT uses always first return value when used AnyTimes()

**Actual behavior**
When I am using the following line multiple times with different return objects then the mocked function returns always the object like expected:

```go
mock.EXPECT().DoSomething(gomock.Any(), gomock.Any()).Times().Return(someObject1, nil)
// DoSomething() returns the object someObject1
mock.EXPECT().DoSomething(gomock.Any(), gomock.Any()).Times().Return(someObject2, nil)
// DoSomething() returns the object someObject2
```

But if I change the Times() call to AnyTimes(), then it returns always the first given return object:

```go
mock.EXPECT().DoSomething(gomock.Any(), gomock.Any()).AnyTimes().Return(someObject1, nil)
// DoSomething() returns the object someObject1
mock.EXPECT().DoSomething(gomock.Any(), gomock.Any()).AnyTimes().Return(someObject2, nil)
// DoSomething() returns the object someObject1
```

**Expected behavior**
I expected, that it should return the new set return object like it is doing on the first example with Time().

- 原因

because AnyTimes() is greedy. There is no reason for it to move on to the second mock because nothing new has changed and the first EXPECT is still doing its work.

相关 Issue: https://github.com/golang/mock/issues/402

## 参考链接

- https://github.com/golang/mock

