[TOC]

## go monkey 介绍

用来给变量、函数和方法打桩。注意：**moneky 不是线程安全的，不能用在并发测试中**。

- 安装库

~~`github.com/agiledragon/gomonkey`~~ 旧版本不再建议使用，因为不支持 private 方法打桩。

`github.com/agiledragon/gomonkey/v2` 新版本 support a patch for a private member method。使用方法 `ApplyPrivateMethod`。

我这里使用的版本是 `github.com/agiledragon/gomonkey/v2 v2.7.0`

- Test Method

```bash
$ cd test 
$ go test -gcflags=all=-l
```

**注意：**`-gcflags=all=-l` 这一行命令表示关闭编译器内联优化。不关闭可能导致 mock 失败。

## go monkey 使用

### mock 函数

- ApplyFunc

假设我们要调用远程服务器的某个方法。在本地跑测试用例的时候我们不可能真正的建立远程连接，所以需要进行打桩。

```go
package main

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

// 远程服务器方法
func callRemoteAdd(a, b int) (int, error) {
	// do something in remote computer
	return 0, nil
}

// 本地方法
func compute(a, b int) (int, error) {
	sum, err := callRemoteAdd(a, b)
	return sum, err
}

func TestCompute(t *testing.T) {
    // 对远程服务器的方法进行打桩
	patches := gomonkey.ApplyFunc(callRemoteAdd, func(a, b int) (int, error) {
        return 2, nil
	})
	defer patches.Reset() // 使用 patches.Reset() 来恢复这次打桩前的状态

	sum, err := compute(1, 1)
	if sum != 2 || err != nil {
		t.Errorf("expected %v, got %v", 2, sum)
	}
}

```

### mock 方法

- ApplyMethod

```go
type computer struct {
}

func (c *computer) CallRemoteAdd(a, b int) (int, error) {
	// do something in remote computer
	return -1, nil
}

func (c *computer) compute(a, b int) (int, error) {
	sum, err := c.CallRemoteAdd(a, b)
	return sum, err
}

func TestComputeMethod(t *testing.T) {
	var c *computer
	patches :=gomonkey.ApplyMethod(reflect.TypeOf(c), "CallRemoteAdd", func(_ *computer, a, b int) (int, error) {
		return 2, nil
	})
	defer patches.Reset()

	sum, err := c.compute(1, 1)
	if sum != 2 || err != nil {
		t.Errorf("expected %v, got %v", 2, sum)
	}
}
```

- 打桩 private 方法: `ApplyPrivateMethod`

### mock 接口

参考（一定要先看这里的代码，fake 包中的实现类和接口的定义）：https://github.com/agiledragon/gomonkey/blob/master/test/apply_interface_reused_test.go

当我们定义了一个 interface 时，系统中就会有一个或者多个实现类(struct)，我们可以通过 ApplyFunc 接口让 interface 变量指向一个实现类对象，然后通过 ApplyMethod 接口来改变该实现类的行为，这就相当于对 interface 完成了打桩。

示例代码：先构造一个 Etcd 对象 e，通过第一层 convey 调用 ApplyFunc 让 Db 的 interface 变量指向 e，然后在第二层 convey 中调用 ApplyMethod 对 Db 完成打一个桩。

注意：**不同包**下对接口进行打桩，需要任一实现类(struct)为导出（首字母大写）才可以实现接口打桩。**以下代码中，如果 fake 包中的 Etcd 和 Mysql 类定义为非导出，那么 ApplyMethod 是无法实现的。因为 ApplyMethod 的参数需要具体的实现类。**

```go
func TestApplyInterfaceReused(t *testing.T) {
    e := &fake.Etcd{}

    Convey("TestApplyInterface", t, func() {
        patches := ApplyFunc(fake.NewDb, func(_ string) fake.Db {
            return e
        })
        defer patches.Reset()
        db := fake.NewDb("mysql") // new 一个 mysql 的 db。这里因为上面打桩过，所以这里其实返回的实例是 Etcd。

        Convey("TestApplyInterface", func() {
            info := "hello interface"
            // 因为 Mysql 和 Etcd 都是 db 的实现类。所以对 Etcd 的 Retrieve 方法进行打桩，就相当于对 db 接口方法 Retrieve 打桩。
            patches.ApplyMethod(e, "Retrieve",
                func(_ *fake.Etcd, _ string) (string, error) {
                    return info, nil
                })
            output, err := db.Retrieve("") // 调用 db 接口的 Retrieve 方法
            So(err, ShouldEqual, nil)
            So(output, ShouldEqual, info)
        })
    })
}
```

### mock 全局变量

- ApplyGlobalVar

```go
var num = 10

func TestGlobalVar(t *testing.T) {
	patches := gomonkey.ApplyGlobalVar(&num, 12)
	defer patches.Reset()

	if num != 12 {
		t.Errorf("expected %v, got %v", 12, num)
	}
}
```

### mock 函数序列

- ApplyFuncSeq

方法序列也同理

```go
func getRandomString() string {
	return "hello world"
}

func TestGetRandomString(t *testing.T) {
	outputs := []gomonkey.OutputCell{
		{Values: gomonkey.Params{"123"}, Times: 2}, // 模拟函数的第 1,2 次输出。Times 表示输出的次数。
		{Values: gomonkey.Params{"456"}},           // 模拟函数的第 3 次输出
		{Values: gomonkey.Params{"789"}},           // 模拟函数的第 4 次输出
	}

	// 打桩
	patches := gomonkey.ApplyFuncSeq(getRandomString, outputs)
	defer patches.Reset()

	fmt.Printf("getRandomString(): %v\n", getRandomString())
	fmt.Printf("getRandomString(): %v\n", getRandomString())
	fmt.Printf("getRandomString(): %v\n", getRandomString())
	fmt.Printf("getRandomString(): %v\n", getRandomString())
}
```

OutputCell 的 Values 代表返回的值。Times 表示返回几次。

```go
type Params []interface{}
type OutputCell struct {
	Values Params
	Times  int
}
```

运行结果输出

```go
=== RUN   TestGetRandomString
getRandomString(): 123
getRandomString(): 123
getRandomString(): 456
getRandomString(): 789
--- PASS: TestGetRandomString (0.00s)
```

## 新版本特性

### 支持打桩 private 方法

[gomonkey支持为private method打桩了](https://www.jianshu.com/p/7546e788613b)

### gomonkey 的惯用法更新

> 由于当前写文档的时候 gomonkey 还是旧版本，导致新版本部分更新的内容无法体现文档中，后续如果有空我再更新一遍文档

后续作者优化了 gomonkey 打桩的方法，使的打桩更加方便简洁：[你该刷新gomonkey的惯用法了](https://www.jianshu.com/p/25d49af216b7?utm_campaign=hugo&utm_medium=reader_share&utm_content=note&utm_source=weixin-friends)

## gomonkey 踩坑

### ~~private 方法无法打桩~~

~~Monkey 框架的实现中大量使用了反射机制，但是，go1.6 版本和更高版本（比如go1.7）的反射机制有些差异：~~

> ~~在 **go1.6 版本中反射机制会导出所有方法**（不论首字母是大写还是小写），而在**更高版本中反射机制仅会导出首字母大写的方法**。~~

~~反射机制的这种差异导致了 Monkey 框架的第二个缺陷：在 go1.6 版本中可以成功打桩的首字母小写的方法，当 go 版本升级后 Monkey 框架会显式触发 panic。~~

### 记得 patches.Reset

没有执行到 patches.Reset() 之前，**mock 方法都会一直生效**，可能会导致影响到其他的测试用例。

```go
package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

type computer struct {
}

func (c *computer) CallRemoteAdd(a, b int) (int, error) {
	// do something in remote computer
	return -1, nil
}

func (c *computer) compute(a, b int) (int, error) {
	sum, err := c.CallRemoteAdd(a, b)
	return sum, err
}

func TestComputeMethod(t *testing.T) {
	var c *computer
	patches := gomonkey.ApplyMethod(reflect.TypeOf(c), "CallRemoteAdd", func(_ *computer, a, b int) (int, error) {
		return 2, nil // 打桩返回 2
	})
    // 如果不加 patches.Reset() 那么 TestComputeMethod01 也会受到影响。两次用例输出都是 sum: 2
    // 如果添加了 patches.Reset()，那么 TestComputeMethod 用例输出是 sum: 2 
    // TestComputeMethod01 用例输出是 sum: -1
    defer patches.Reset() 
    
	test01()
}

func TestComputeMethod01(t *testing.T) {
	test01()
}

func test01() {
	var c *computer
	sum, _ := c.compute(1000, 1000)
	fmt.Printf("sum: %v\n", sum)
}

```

## 参考链接

- [go 单元测试 gomonkey](https://www.cnblogs.com/lanyangsh/p/14587921.html)
- https://pkg.go.dev/github.com/agiledragon/gomonkey/v2@v2.7.0#section-readme
- [gomonkey支持为private method打桩了](https://www.jianshu.com/p/7546e788613b)
- [你该刷新gomonkey的惯用法了](https://www.jianshu.com/p/25d49af216b7?utm_campaign=hugo&utm_medium=reader_share&utm_content=note&utm_source=weixin-friends)