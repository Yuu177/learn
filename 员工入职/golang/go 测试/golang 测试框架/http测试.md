[TOC]

# http 测试

## 创建真实网络连接测试 http handler 接口

假设需要**测试某个 http API 接口**的 handler 能够正常工作，例如 `helloHandler`

```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
```

那我们可以创建真实的网络连接进行测试：

```go
// test code
import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

func handleError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("failed", err)
	}
}

func TestConn(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	handleError(t, err)
	defer ln.Close()

	http.HandleFunc("/hello", helloHandler)
	go http.Serve(ln, nil)

	resp, err := http.Get("http://" + ln.Addr().String() + "/hello")
	handleError(t, err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	handleError(t, err)

	if string(body) != "hello world" {
		t.Fatal("expected hello world, but got", string(body))
	}
}
```

- `net.Listen("tcp", "127.0.0.1:0")`：监听一个未被占用的端口，并返回 Listener。
- 调用 `http.Serve(ln, nil)` 启动 http 服务。
- 使用 `http.Get` 发起一个 Get 请求，检查返回值是否正确。

## httptest

### 作为 server 测试 http handler 接口

针对 http api 开发的场景，使用**标准库 net/http/httptest** 进行测试更为高效。

上述的测试用例改写如下：

```go
// test code
import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConn(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	helloHandler(w, req)
	bytes, _ := ioutil.ReadAll(w.Result().Body) // 回放 ResponseRecorder 记录接口处理后响应的内容。

	if string(bytes) != "hello world" {
		t.Fatal("expected hello world, but got", string(bytes))
	}
}
```

使用 httptest 模拟请求对象(req)和响应对象(w)，达到了相同的目的。

#### 使用 httptest 对 gin 框架测试

HTTP 测试首选 net/http/httptest 包。

```go
package main

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    return r
}

func main() {
    r := setupRouter()
    r.Run(":8080")
}
```

上面这段代码的测试用例：

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    router.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)
    assert.Equal(t, "pong", w.Body.String())
}
```

### 作为 client 发送 http 请求

当我们作为 client 方要调用其他服务的 http 接口的时候，for example：

- 业务代码

```go
package user

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *client {
	return &client{baseURL: baseURL, httpClient: http.DefaultClient}
}

// 推荐写法，可以通过 httptest 来测试
func (c *client) GetInfo() (string, error) {
	url := fmt.Sprintf("%s/info", c.baseURL)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Printf("http do error|err:%v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}

// 不建议的写法
// 无法通过 httptest 来 mock，因为该请求的 base URL 无法请求到我们的 httptest server 上来。
// 如果非要这么写，可以通过 httpmock 这个框架来测试。下面会介绍到。
func GetInfo() (string, error) {
	url := fmt.Sprintf("%s/info", "https://test.com")
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("http do error|err:%v\n", err)
		return "", err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}
```

- 测试代码

```go
package user_test

import (
	user "demo/http_unit"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
    "github.com/stretchr/testify/assert"
)

func TestUserClientGetInfo(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'\n", r.Method)
		}
		if r.URL.EscapedPath() != "/info" { // pattern
			t.Errorf("Expected request to '/info', got '%s'\n", r.URL.EscapedPath())
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("httptest"))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(handler)) // 所有访问这个 httptest server 的请求都由这个 handler 函数去处理
	defer fakeServer.Close() // 别忘了 close
	
	userClient := user.NewClient(fakeServer.URL) // 使用 fakeServer 的 base URL。base URL of form http://ipaddr:port with no trailing slash
	info, err := userClient.GetInfo()
	if err != nil {
		t.Errorf("get info error|err:%v\n", err)
	}
	assert.Equal(t, "httptest", info)
}
```

1、启动一个 httptest server，并编写 handler 方法

2、客户端访问 httptest server 的 URL

## httpmock

Easy mocking of http responses from external resources.

- 一个简单的 demo

```go
package user_test

import (
	user "demo/http_unit"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
    "github.com/stretchr/testify/assert"
)

func TestUserClientGetInfo(t *testing.T) {
	// ...
}

func TestGetInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset() // 别忘了 DeactivateAndReset

	httpmock.RegisterResponder("GET", "https://test.com/info", httpmock.NewStringResponder(200, string("httpmock")))

	info, err := user.GetInfo()
	if err != nil {
		t.Errorf("get info error|err:%v\n", err)
	}
	assert.Equal(t, "httpmock", info)
}
```

1、调用 `Activate` 方法启动 httpmock 环境。

2、在 defer 里面调用 `DeactivateAndReset` 结束 mock。

3、通过 `httpmock.RegisterResponder` 方法进行 mock 规则注册。

4、这时候再通过 http client 发起的请求就都会被 httpmock 拦截，如果匹配到刚刚注册的规则就会按照注册的内容返回对应 response。如果匹配不到，显示 `no responder found`。

5、规则注册还可以通过正则表达式来匹配

## 参考链接

- [Golang单元测试实战一：httpMock](https://www.jianshu.com/p/545963b593de)

- [Gin 框架中文文档：怎样编写 Gin 的测试用例](https://learnku.com/docs/gin-gonic/1.7/testing/11359)

- https://github.com/jarcoal/httpmock