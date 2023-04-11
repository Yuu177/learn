[TOC]

# Kafka mock

## Shopify sarama

项目中，大部分用到的是这个 kafka 库 `github.com/Shopify/sarama`。

这个库提供了  mock 代码(`github.com/Shopify/sarama/mocks`)来方便我们对 kafka 进行测试。

- 代码 example：https://github.com/Shopify/sarama/tree/main/mocks

>下面内容是对 example 的补充

### check the message 方法

> 检验 kafka 消息：包括 topic，key，value 等

```go
type MessageChecker func(*sarama.ProducerMessage) error

func (mp *AsyncProducer) ExpectInputWithMessageCheckerFunctionAndSucceed(cf MessageChecker) *AsyncProducer
```

### check the message value 方法

> 检验 kafka 消息的 value

```go
type ValueChecker func(val []byte) error

func (mp *AsyncProducer) ExpectInputWithCheckerFunctionAndSucceed(cf ValueChecker) *AsyncProducer
```

### 代码示例

```go
package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
)

func TestProducerWithBrokenPartitioner(t *testing.T) {
	mp := mocks.NewAsyncProducer(t, nil)
	mp.ExpectInputWithMessageCheckerFunctionAndSucceed(func(msg *sarama.ProducerMessage) error { // check the message
		if msg.Topic != "test" {
			t.Errorf(`Expected topic "test", found: %v\n`, msg.Topic)
		}

		gotValue, _ := msg.Value.Encode()
		expValue, _ := sarama.ByteEncoder("messageCheckerFunc").Encode()
		if !reflect.DeepEqual(gotValue, expValue) {
			t.Errorf("gotValue: %v, expValue: %v\n", string(gotValue), string(expValue))
		}
		return nil
	})

	mp.ExpectInputWithCheckerFunctionAndSucceed(func(val []byte) error { // check the message value
		re := "checkerFunc"
		matched, err := regexp.MatchString(re, string(val))
		if err != nil {
			return errors.New("Error while trying to match the input message with the expected pattern: " + err.Error())
		}
		if !matched {
			return fmt.Errorf("No match between input value \"%s\" and expected pattern \"%s\"", val, re)
		}
		return nil
	})

	msg := &sarama.ProducerMessage{
		Topic:     "test",
		Partition: 1,
		Key:       sarama.StringEncoder("user"),
		Value:     sarama.ByteEncoder("messageCheckerFunc"),
	}
	mp.Input() <- msg

	msg1 := &sarama.ProducerMessage{
		Topic: "test",
		Value: sarama.ByteEncoder("checkerFunc"),
	}
	mp.Input() <- msg1

	if err := mp.Close(); err != nil { // 需要关闭才会发送
		t.Error(err)
	}
}
```

