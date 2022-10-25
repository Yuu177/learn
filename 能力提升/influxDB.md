[TOC]

# influxDB

## 常用命令

```sql
# 显示所有的数据库
show databases
# 使用数据库 
use db_name
# 显示该数据库中所有的表
show measurements
# influx 默认是 纳米保存时间戳，可以修改为 UTC 时间 
precision rfc3339
# 修改为 UTC 时间显示后，在查询语句后面加上 tz('Asia/Shanghai') 改为东八区时间显示
select * from measurement_name tz('Asia/Shanghai')
# 删除表数据 
delete from measurement_name
# 如何知道哪些字段是tags，哪些字段是 fields
show field keys from measurement;
show tag keys from measurement;
# distinct 去重
select distinct(uid) from usertable
```

influxdb 没有写入去重，但是可以通过 distinct 查询去重。

distinct 的字段仅是 field 不能是 tag。group by 仅是 tag 不能是 field。

同时有 group by 和 distinct 时，先 group by 后 distinct。使用子查询 select * from (select distinct(msg_id) from whisper) group by time(10s)

## 数据丢失问题（influxDB 缺陷）

做的一个埋点的需求，然后客户端上报数据吗，我们服务端插入数据库的时候 point 丢失。

**A point is uniquely identified by the measurement name, tag set, and timestamp**

因为我们服务端向数据库插入的时候，是一个 for 循环，插入的时间间隔太小的了（timestamp 的精度不够高）。所以导致会插入相同的 timestamp。加上 points 的 measurement 和 tag 都一样，导致了前面插入的 points 被覆盖掉了。

解决办法：

设置一个唯一的 tag。（一个字段只能设置为 tag 或者 filed）

参考链接：https://juejin.cn/post/6844903887762112525