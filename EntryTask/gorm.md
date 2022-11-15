[TOC]

## 使用 gorm 查询数据

执行查询的函数，gorm 提供下面几个查询函数：

- **Take 查询一条记录**

```sql
// 定义接收查询结果的结构体变量
food := Food{} 
// 等价于：SELECT * FROM `foods` LIMIT 1
db.Take(&food)
```

- **First 查询一条记录，根据主键 ID 排序(正序)，返回第一条记录**

```sql
//等价于：SELECT * FROM `foods`  ORDER BY `foods`.`id` ASC LIMIT 1    
db.First(&food)
```

- **Last 查询一条记录, 根据主键ID排序(倒序)，返回第一条记录**

```sql
// 等价于：SELECT * FROM `foods` ORDER BY `foods`.`id` DESC LIMIT 1  
// 语义上相当于返回最后一条记录 
db.Last(&food)
```

- **Find 查询多条记录，Find函数返回的是一个数组**

```sql
// 因为 Find 返回的是数组，所以定义一个商品数组用来接收结果
var foods []Food 
// 等价于：SELECT * FROM `foods` 
db.Find(&foods)
```

- **Pluck 查询一列值**

```sql
// 商品标题数组
var titles []string 
// 返回所有商品标题 
// 等价于：SELECT title FROM `foods` 
// Pluck提取了title字段，保存到titles变量 
// 这里Model函数是为了绑定一个模型实例，可以从里面提取表名。 
db.Model(&Food{}).Pluck("title", &titles)
```

## 示例代码

```go
package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// UserInfo 用户信息
type UserInfo struct {
	ID     uint
	Name   string
	Gender string
	Hobby  string
}

func main() {
	db, err := gorm.Open("mysql", "user:user@(127.0.0.1:3306)/testdb01?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 自动迁移。这个时候 UserInfo 的数据库表会自动建立
	db.AutoMigrate(&UserInfo{})

	u1 := UserInfo{1, "枯藤", "男", "篮球"}
	u2 := UserInfo{2, "topgoer.com", "女", "足球"}
	// 创建记录（插入）
	db.Create(&u1)
	db.Create(&u2)

	// 查询
	var u = new(UserInfo)
	db.First(u) // 查询出来的数据保存到 u 中。
	fmt.Printf("%#v\n", u)

	var uu UserInfo
	db.Find(&uu, "hobby=?", "足球")
	fmt.Printf("%#v\n", uu)

	// 更新
	db.Model(&u).Update("hobby", "双色球") // 只更新 hobby 字段
	// 更新整条记录
	u.Gender = "女"
	u.Hobby = "football"
	db.Save(&u)

	// 删除
	db.Delete(&u)
}
```

## 参考文档

https://gorm.io/docs/