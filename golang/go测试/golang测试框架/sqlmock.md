[TOC]

# sqlmock

另一个项目中比较常见的依赖其实就是数据库，在遇到数据库的依赖时，我们一般都会使用 [sqlmock](https://github.com/DATA-DOG/go-sqlmock) 来模拟数据库的连接。

## 正则匹配

sqlmock 中默认使用的是正则表达式去匹配 sql 语句。

如 `sqlmock.ExpectQuery()` 和 `sqlmock.ExpectExec()` 等，如果想直接匹配 sql 语句需要加上 `regexp.QuoteMeta()`

- example

```go
result := sqlmock.NewRows([]string{"user_id"}).AddRow(777)
// 直接匹配 sql 语句
mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `my_tab` WHERE (enabled=1 AND create_time BETWEEN 0 AND 0)")).WillReturnRows(result)
// 使用正则来匹配 sql 语句，需要使用 \\ 来转义字符
mock.ExpectQuery("SELECT \\* FROM `my_tab` WHERE \\(enabled=1 AND create_time BETWEEN 0 AND 0\\)").WillReturnRows(result)
```

## 参考链接

- https://github.com/DATA-DOG/go-sqlmock