[TOC]

# golang 单元测试

## 如何写好单元测试

单元测试的好处：略

函数/方法写法不同，测试难度也是不一样的。职责单一，参数类型简单，与其他函数耦合度低的函数往往更容易测试。

『这种代码没法测』这种时候，就得思考函数的写法可不可以改得更好一些。为了代码可测试而重构是值得的。（~~现实中是不可能的，特别祖传代码还是不要动好~~） 

下面介绍如何使用 Go 语言的标准库 testing 进行单元测试。

## 测试函数简单介绍

go test 命令是一个按照一定的约定和组织来测试代码的程序。在包目录内，所有以 `_test.go` 为后缀名的源文件在执行 go build 时不会被构建成包的一部分，它们是 go test 测试的一部分。

在 `*_test.go` 文件中，有三种类型的函数：**单元测试函数、基准测试(benchmark)函数**、示例函数（这个就不介绍了）。

## 单元测试函数

### 单元测试函数介绍

Go 语言推荐测试文件和源代码文件放在一块，测试文件以 _test.go 结尾。

比如，当前文件夹下有 calc.go 一个文件，我们想测试 calc.go 中的 `Add` 和 `Mul` 函数，那么应该新建 calc_test.go 作为测试文件。

```bash
example/
   |--calc.go
   |--calc_test.go
```

calc.go 的代码如下：

```go
package main

func Add(a int, b int) int {
    return a + b
}

func Mul(a int, b int) int {
    return a * b
}
```

那么 calc_test.go 中的测试用例可以这么写：

```go
package main

import "testing"

func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		t.Errorf("1 + 2 expected be 3, but %d got", ans)
	}

	if ans := Add(-10, -20); ans != -30 {
		t.Errorf("-10 + -20 expected be -30, but %d got", ans)
	}
}
```

**运行测试用例** `go test -v -cover`

- go test，该文件夹下所有的测试用例都会被执行。

- -v 参数会显示每个用例的测试结果。
- -cover 参数可以查看覆盖率（分支覆盖率）。

**结果输出**

```bash
=== RUN   TestAdd
--- PASS: TestAdd (0.00s)
PASS
coverage: 100.0% of statements
ok      tpy/example     0.606s
```

如果只想运行其中的一个用例，例如 TestAdd，可以用 `-run` 参数指定，该参数支持通配符 `*`，和部分正则表达式。`go test -run TestAdd`

**小结**

- 测试文件以 `_test.go` 结尾，通常建议和源文件同一个目录。

- 测试用例名称一般命名为 Test 加上待测试的方法名（驼峰或者下划线）。且参数有且只有一个，在这里是 `t *testing.T`。
- `go test` 运行该文件夹下所有的测试用例。

### 子测试(Subtests)

子测试是指我们可以在单元测试中启动多个测试用例。

子测试是 Go 语言内置支持的，可以在某个测试用例中，根据测试场景使用 `t.Run` 创建不同的子测试用例：

```go
// calc_test.go
func TestMul(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail")
		}

	})
	t.Run("second", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}
```

- 之前的例子测试失败时使用 `t.Error/t.Errorf`，这个例子中使用 `t.Fatal/t.Fatalf`，区别在于前者遇错不停，还会继续执行其他的测试用例，后者遇错即停。

运行测试用例的子测试：

```go
$ go test -run TestMul/first -v
=== RUN   TestMul
=== RUN   TestMul/first
--- PASS: TestMul (0.00s)
    --- PASS: TestMul/first (0.00s)
PASS
ok      example 0.008s
```

对于多个子测试的场景，更推荐如下的写法，表格驱动(table-driven tests)

```go
//  calc_test.go
func TestMul(t *testing.T) {
	cases := []struct {
		Name           string
		A, B, Expected int
	}{
		{"pos", 2, 3, 6},
		{"neg", 2, -3, -6},
		{"zero", 2, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.A, c.B); ans != c.Expected {
				t.Fatalf("%d * %d expected %d, but %d got",
					c.A, c.B, c.Expected, ans)
			}
		})
	}
}
```

所有用例的数据组织在切片 cases 中，看起来就像一张表，借助循环创建子测试。这样写的好处有：

- 新增用例非常简单，只需给 cases 新增一条测试数据即可。
- 测试代码可读性好，直观地能够看到每个子测试的参数和期待的返回值。
- 用例失败时，报错信息的格式比较统一，测试报告易于阅读。

如果数据量较大，或是一些二进制数据，推荐使用相对路径从文件中读取。

#### 子测试模板

> 推荐下面这个多个子测试模板写法

```go
func Cal(a int, b int) (int, int) {
	return 0, 0
}

// test func
func TestCal(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name  string // 子测试名称
		args  args   // 函数入参
		want  int    // 期望值
		want1 int    // 期望值 1
	}{
		// Add test cases.
		{"first", args{1, 2}, 0, 0}, // 测试用例只需要在这里维护即可
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Cal(tt.args.a, tt.args.b)
			if got != tt.want {
				t.Errorf("Cal() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cal() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
```

**小结**

- `t.Run` 创建子测试用例
- `go test -run TestCases/subCase` 运行某个用例下的某个子测试用例
- 多个子测试的场景建议采用 table-driven tests 写法

### 帮助函数(helpers)

对一些重复的逻辑，抽取出来作为公共的帮助函数(helpers)，可以增加测试代码的可读性和可维护性。 借助帮助函数，可以让测试用例的主逻辑看起来更清晰。

例如，我们可以将创建子测试的逻辑抽取出来：

```go
// calc_test.go
package main
 
import "testing"
 
type calcCase struct{ A, B, Expected int }
 
func createMulTestCase(t *testing.T, c *calcCase) {
    // t.Helper() // 注释掉 t.Helper() 对比前后运行输出
    if ans := Mul(c.A, c.B); ans != c.Expected {
        t.Fatalf("%d * %d expected %d, but %d got",
            c.A, c.B, c.Expected, ans)
    }
}
 
func TestMul(t *testing.T) {
    createMulTestCase(t, &calcCase{2, 3, 6})
    createMulTestCase(t, &calcCase{2, -3, -6})
    createMulTestCase(t, &calcCase{2, 0, 1}) // wrong case
}
```

在这里，我们故意创建了一个错误的测试用例，运行 go test，用例失败，会报告错误发生的文件和行号信息：

```bash
$ go test
--- FAIL: TestMul (0.00s)
    calc_test.go:11: 2 * 0 expected 1, but 0 got
FAIL
exit status 1
FAIL    example 0.007s
```

可以看到，错误发生在第 11 行，也就是帮助函数 `createMulTestCase` 内部。但 18, 19, 20 行都调用了该方法，我们第一时间并不能够确定是哪一行发生了错误。报错信息都在同一处，不方便问题定位。

因此，Go 语言在 1.9 版本中引入了 `t.Helper()`，用于标注该函数是帮助函数，报错时将输出帮助函数调用者的信息，而不是帮助函数的内部信息。

修改 `createMulTestCase`，取消注释 `t.Helper()`

运行 go test，报错信息如下，可以非常清晰地知道，错误发生在第 20 行。

```bash
$ go test
--- FAIL: TestMul (0.00s)
    calc_test.go:20: 2 * 0 expected 1, but 0 got
FAIL
exit status 1
FAIL    example 0.006s
```

关于 helper 函数的 建议：

- 不要返回错误， 帮助函数内部直接使用 `t.Error` 或 `t.Fatal` 即可，在用例主逻辑中不会因为太多的错误处理代码，影响可读性。

**小结**

- 对一些重复的逻辑，抽取出来作为公共的帮助函数，增加测试代码的可读性和可维护性。
- 帮助函数中调用 `t.Helper()` 让报错信息更准确，有助于定位。

### setup 和 teardown

> 不太常用

如果在同一个测试文件中，每一个测试用例运行前后的逻辑是相同的，一般会写在 `setup` 和 `teardown` 函数中。例如执行前需要实例化待测试的对象，如果这个对象比较复杂，很适合将这一部分逻辑提取出来；

执行后，可能会做一些资源回收类的工作，例如关闭网络连接，释放文件等。标准库 testing 提供了这样的机制：

```go
func setup() {
	fmt.Println("Before all tests")
}

func teardown() {
	fmt.Println("After all tests")
}

func Test1(t *testing.T) {
	fmt.Println("I'm test1")
}

func Test2(t *testing.T) {
	fmt.Println("I'm test2")
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
```

- 在这个测试文件中，包含有 2 个测试用例，Test1 和 Test2。
- 如果测试文件中包含函数 `TestMain(m *testing.M)`，那么生成的测试将调用 `TestMain(m *testing.M)`，而不是直接运行测试。
- 调用 `m.Run()` 触发所有测试用例的执行，并使用 `os.Exit()` 处理返回的状态码，如果不为 0，说明有用例失败。
- 因此可以在调用 `m.Run()` 前后做一些额外的准备(setup)和回收(teardown)工作。

执行 go test，将会输出

```bash
$ go test
Before all tests
I'm test1
I'm test2
PASS
After all tests
ok      example 0.006s
```

**小结**

- 如果在同一个测试文件中，每一个测试用例运行前后的逻辑是相同的，建议使用 `TestMain(m *testing.M)`。

### 以并行的方式进行测试

单元测试的目的是独立地进行测试。尽管有些时候，测试套件会因为内部存在依赖关系而无法独立地进行单元测试，但是只要单元测试可以**独立地进行**，用户就可以通过并行地运行测试用例来提升测试的速度了。

```go
package main

import (
	"testing"
	"time"
)

func TestParallel1(t *testing.T) {
	t.Parallel()
	time.Sleep(1 * time.Second)
}

func TestParallel2(t *testing.T) {
	t.Parallel()
	time.Sleep(2 * time.Second)
}

func TestParallel3(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}
```

结果输出

- 未加 `t.Parallel()`

```bash
tpy@C02G65GVMD6M:~/codeTest/goLangTest/unitTesting/parallelTest % go test -v
=== RUN   TestParallel1
--- PASS: TestParallel1 (1.00s)
=== RUN   TestParallel2
--- PASS: TestParallel2 (2.00s)
=== RUN   TestParallel3
--- PASS: TestParallel3 (3.00s)
PASS
ok      tpy/parallelTest        6.340s
```

- 添加 `t.Parallel()` 代码后

```bash
tpy@C02G65GVMD6M:~/codeTest/goLangTest/unitTesting/parallelTest % go test -v
=== RUN   TestParallel1
=== PAUSE TestParallel1
=== RUN   TestParallel2
=== PAUSE TestParallel2
=== RUN   TestParallel3
=== PAUSE TestParallel3
=== CONT  TestParallel1
=== CONT  TestParallel3
=== CONT  TestParallel2
--- PASS: TestParallel1 (1.00s)
--- PASS: TestParallel2 (2.00s)
--- PASS: TestParallel3 (3.00s)
PASS
ok      tpy/parallelTest        3.336s
```

**小结**

只要单元测试可以独立地进行，用户就可以通过 `t.Parallel()` 并行地运行测试用例来提升测试的速度。

## 基准测试(Benchmark)

### 基准测试介绍

go 中基准测试非常直观：测试程序要做的就是将被测试的代码执行了 `b.N` 次，以便准确地检测出代码的响应时间，其中 `b.N` 的值将根据被执行的代码而改变。

基准测试用例的定义如下：

```go
func BenchmarkName(b *testing.B){
    for i := 0; i < b.N; i++ {
		// code
	}
}
```

- 函数名必须以 Benchmark 开头，后面一般跟待测试的函数名
- 参数为 `b *testing.B`。

例如：

```go
func BenchmarkHello(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
```

运行 `go test -bench . -benchmem`

- 执行基准测试时，需要添加 `-bench` 参数，`.` 表示匹配该文件夹下所有所有的基准测试。可以替换为正则表达式或者某个指定的基准测试。
- `-benchmem` 可以提供每次操作分配内存的次数，以及每次操作分配的字节数。 

```bash
$ go test -bench . -benchmem
goos: darwin
goarch: amd64
pkg: tpy/bechmarkTest
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkHello-12       15309099                65.35 ns/op           16 B/op          1 allocs/op
PASS
ok      tpy/bechmarkTest        1.438s
```

对第 6 行的结果输出介绍

| GOMAXPROCS | 运行了多少次 | 平均每次耗时 | 每次操作分配的字节数 | 每次操作进行内存的分配次数 |
| ---------- | ------------ | ------------ | -------------------- | -------------------------- |
| -12        | 15309099     | 65.35 ns/op  | 16 B/op              | 1 allocs/op                |

**注意：**测试代码需要保证函数可重入性及无状态，也就是说，测试代码不使用全局变量等带有记忆性质的数据结构。避免多次运行同一段代码时的环境不一致。

- 如果在运行前基准测试需要一些耗时的配置，则可以使用 `b.ResetTimer()` 先重置定时器，这样可以避免 for 循环之前的初始化代码的干扰。

```go
func BenchmarkHello(b *testing.B) {
    ... // 耗时操作
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
```

**小结**

- 基准测试函数名为 `Benchmark` 开头，参数为 `b *testing.B`。

- `go test -bench .` 执行基准测试，`-benchmem` 参数增加内存相关信息打印。
- 在运行前基准测试需要一些耗时的配置，使用 `b.ResetTimer()` 重置定时器。

### 测试并发性能

- 使用 `b.RunParallel` 测试并发性能
- `RunParallel` 并发的执行 benchmark。`RunParallel` 创建多个 goroutine 然后把 `b.N` 个迭代测试分布到这些 goroutine 上。goroutine 的数目默认是 GOMAXPROCS

```go
func BenchmarkParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
            fmt.Sprintf("hello")
		}
	})
}
```

**小结**

- 使用 `b.RunParallel` 测试并发性能。

## 可视化测试覆盖率

前面我们说过 go test -cover 可以显示出该测试的函数的覆盖率，但是没有具体的信息。所以这里生成可视化的 html 文件来查看代码的覆盖率情况。

1. 首先生成测试覆盖率的中间文件：covprofile

`go test -coverprofile=covprofile`

2. 再通过中间文件生成 html 文件：coverage.html。

`go tool cover -html=covprofile -o coverage.html`

浏览器打开 html 文件即可。

## 补充

### assert package

> 使用 assert 包替代 t.Fatal/t.Error

The `assert` package provides some helpful methods that allow you to write better test code in Go.

```go
package yours

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

  // assert equality
  assert.Equal(t, 123, 123, "they should be equal")

  // assert inequality
  assert.NotEqual(t, 123, 456, "they should not be equal")

  // assert for nil (good for errors)
  assert.Nil(t, object)

  // assert for not nil (good when you expect something)
  if assert.NotNil(t, object) {

    // now we know that object isn't nil, we are safe to make
    // further assertions without causing any errors
    assert.Equal(t, "Something", object.Value)

  }

}
```

### 规范测试文件中 package 名称

```
.
└── mymath
    ├── math.go
    └── math_test.go
```

业务代码 math.go

```go
package mymath

func Add(a int, b int) int {
	return a + b
}
```

测试代码 math_test.go

```go
// test 文件在包名后面加上 _test（只能加 _test，其他的后缀会报错 MismatchedPkgName）
// 因为同一个文件夹下的代码文件只能有一个包名
// 只有 _test.go 文件 package 才能以 _test 结尾来命名
// 这样子更好的把业务代码和测试代码隔离开来
package mymath_test

func TestAdd(t *testing.T) {
	// code
}
```

## 注意事项

**测试用例应该互不影响**，如果执行某个测试用例后，某个全局变量或者带记忆性的数据结构被修改，执行结束后需要把状态进行恢复，避免影响到其他的测试用例。

