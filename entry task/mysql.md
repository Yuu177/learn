[TOC]

## 常用的命令

```sql
# 登陆命令
mysql -h 主机地址 -u 用户名 -p 用户密码
# 用户名和密码都是 user
# 登陆
mysql -u username -p
# 查看所有数据库
show databases;
# 选择数据库
use dbname;
# 查看表结构
desc tablename;
# 查看最大连接数
show variables like 'max_connections';
# 设置最大连接数为 2000（mysql 重新启动后需要重新设置）
set global max_connections=2000;
```

## 往数据库插入 10 000 000 条数据

- 利用 mysql 内存表插入速度快的特点，先存储过程在内存表中生成数据，然后再从内存表插入普通表中。

- 使用批量插入。

### 使用内存表插入

1、创建 mysql 内存表

把普通表创建语句的 `ENGINE=xx` 改为 `ENGINE=MEMORY ` 表明这个是内存表。内存表的速度是要比写在磁盘数据的表的速度是要快很多的。

```sql
CREATE TABLE `users_mem` (
  `user_name` varchar(255) NOT NULL,
  `password` char(32) NOT NULL,
  PRIMARY KEY (`user_name`)
) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2、创建存储过程

创建叫 `insert_user_mem` 的存储过程。往内存表里插入数据。可以把它看作一个函数，两个入参。s 为起始位置，e 为结束。

```sql
DELIMITER ;;
CREATE PROCEDURE `insert_user_mem`(in s int, in e int)
BEGIN
    declare i int default s;
    WHILE i < e DO
        insert into `user_mem`(`user_name`, `nick_name`, `password`, `pic_name`) values (CONCAT('user', CONV(i,10,16)), CONCAT('user', CONV(i,10,16)), md5(CONCAT('user', CONV(i,10,16))), "");
        set i = i+1;
    end WHILE;
END ;;
DELIMITER;
 
# CONCAT 字符连接
# CONV(i,10,16) 把 i 从十进制转换成十六进制
# md5 使用 md5 加密
```

3、调用存储过程

```sql
CALL insert_user_mem(0, 10000000)
```

4、内存表插入数据到普通表

等待内存表插入完成后，在把内存表数据插到普通表

```sql
insert into user select * from user_mem;
```

- 遇见问题

当调用存储过程一段时间后 `The table 'xxx' is full`。这个报错是因为内存表已经满了。所以需要调整内存表的大小。

```sql
# 查看内存表大小。单位为字节
select @@max_heap_table_size;
# 设置内存表大小为 5G
set @@max_heap_table_size=5073741824; # 5 * 1024 * 1024 * 1024
```

内存表大小重新打开 mysql 后会失效（好像会，忘记了，懒得搞了）

设置为新的内存表后，重新执行存储过程发现还是存了和之前一样多的数据。这个时候把内存表 drop 掉然后重新 create，然后重新执行存储过程就行了。

4 279 938 条数据大概占了 5G 的内存。所以可以把内存表数据分几次插入普通表中。实际结果下来，一千万条数据在普通表中就只占了 1G 左右的磁盘空间。但是为什么相同的内存空间存不了这么多数据呢？

因为在执行插入这些操作的过程，会生成各种日志文件，日志文件占的内存远远大于插入的数据。我们可以使用批量操作来插入。

### 批量插入

批量插入可以减少数据库的 IO 读写次数，所以效率会提升很多。

1、创建表结构

```sql
CREATE TABLE `test_user_profile` (
  `user_name` varchar(255) NOT NULL,
  `nick_name` varchar(255) DEFAULT NULL,
  `pic_name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`user_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
```

2、创建存储过程

```sql
DELIMITER ;;
CREATE PROCEDURE insertTest(in sum INT)
BEGIN
# DECLARE 是定义。SET 是赋值
DECLARE count INT DEFAULT 0;
DECLARE i INT DEFAULT 0;
# 使用 set 定义必须带上 @
set @userName = "";
set @nickName = "";
set @picName = "";
set @exesql = "insert into test_user_profile values ";
set @exedata = "";
set count=0;
set i=0;
while count < sum do
    # concat 后会拼接成：,('user0','user0',''),('user1','user1','')
    set @exedata = concat(@exedata, ",(", concat("'user", count), "',", concat("'user", count), "',", "''", ")");
    #select @exedata; 打印变量
    set count = count+1;
    set i = i+1;    
    # 如果拼接够 1000 的 value 就执行插入
    if i%1000=0
    then
        # SUBSTRING 函数从特定位置开始的字符串返回一个给定长度的子字符串
        # 这里第一位为 "," 所以从第二位开始
        set @exedata = SUBSTRING(@exedata, 2);
        set @exesql = concat(@exesql, @exedata);
                  prepare stmt from @exesql;
        execute stmt;    
        DEALLOCATE prepare stmt;
        set @exedata = "";    
    end if;
end while;
# 把剩下的 value 插入
if length(@exedata)>0 
then
    set @exedata = SUBSTRING(@exedata, 2);
    set @exesql = concat(@exesql, @exedata);
    prepare stmt from @exesql;
    execute stmt;    
    DEALLOCATE prepare stmt;
end if;
END ;;
DELIMITER;
```

3、调用存储过程

```sql
call insertTest(10000000)
```