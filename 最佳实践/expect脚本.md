[TOC]

# Expect 脚本快速入门

Expect 脚本是一种自动化工具，用于与交互式命令行应用程序进行交互。它可以编写脚本来模拟用户输入、捕获输出以及根据特定的条件执行操作。

Expect 脚本的原理基于交互式应用程序的工作方式。当一个应用程序需要用户输入时，它会等待输入并在用户键入后继续执行。Expect 脚本利用这个原理，通过监测应用程序的输出并发送相应的输入，实现自动化的交互过程。

Expect 脚本通常由一系列的 `expect` 和 `send` 命令组成。`expect` 命令用于匹配应用程序的输出，可以使用正则表达式或固定的字符串来进行匹配。一旦匹配成功，接下来的 `send` 命令会发送预定的输入到应用程序。

## Ubuntu 安装 Expect

Expect 脚本文件以 `.sh` 或者 `.expect` 结尾

```shell
sudo apt install tcl tk expect
```

## 常用命令

- `spawn`：启动一个新的交互式进程
- `expect`：从进程接收字符串，期望获得字符串
- `send`：向进程发送字符串，用于模拟用户的输入，注意最后要加 `\r`（回车）或 `\n`（换行）
- `interact`：用户交互。表示执行完成后保持交互状态，把控制权交给控制台
- `sleep n`：使脚本暂停给定的秒数
- `$argc`：脚本的参数个数
- `$argv`：脚本的参数数组，使用 `[lindex $argv n]` 获取第 n 个参数，0 是第一个参数
- `exp_continue`：继续执行后续 `expect`。`exp_continue` 附加于某个 `expect` 判断项之后，可以使该项被匹配后 ，还能继续匹配该 `expect` 判断语句内的其他项
- `expect eof`：结束当前 `spawn` 开启的进程

## expect 命令

> expect 可以接收一个字符串参数，也可以接收正则表达式参数。

- 写法一

```nginx
expect "hello" { send "hello world\r" }
```

只有匹配到 hello，才会输出 hello world。

- 写法二

```nginx
expect "hello"
send "hello world\r"
```

就算匹配不到 hello，一段时间后也会输出 hello world。

### 单一分支匹配模式

单一匹配就是只有一种匹配情况。类似 `if` 语句。

```nginx
expect "hello" { send "hello world\r" }
```

### 多分支匹配模式

类似 `switch case` 语句。

```nginx
expect {
    "hello" { send "hello world\r" }
    "hi" { send "hi world\r" }
    "bye" { send "bye world\r" }
}
```

## 例子一：SSH 登录

```nginx
#!/usr/bin/expect

spawn ssh root@172.16.2.2
expect {
    "yes/no" { send "yes\r"; exp_continue }
    "password:" { send "123\r" }
    "#" { send "echo login ok\r" }
}
expect "#" { send "ls\r" }
expect "#" { send "exit\r" }
expect eof
```

## 例子二：通过跳板机上传版本

我们软件开发的时候，通常需要频繁地更换版本。这个过程可能需要经过多个跳板机，过程非常繁琐且没有意义，这个时候 Expect 脚本就非常好用。

```nginx
#!/usr/bin/expect

# 获取第 1 个参数
set file_name [lindex $argv 0]
set root_path "/home/jinx/Downloads"
set android_path "/data/local/tmp"
set car_ip "161.0.0.2"
set car_path "/var"
set timeout -1

# 本地 push 文件到安卓
spawn adb push $root_path/$file_name $android_path
# 等待文件上传完成
# 在 expect eof 行使用 wait 命令等待前一个 spawn 的进程的结束
expect eof
wait

# 登录到安卓
spawn adb shell
# 从安卓 put 文件到板子上
expect "/ #" { send "busybox ftpput -u root $car_ip $car_path/$file_name $android_path/$file_name\r" }
# 登录板子
expect "/ #" { send "busybox telnet $car_ip\r" }
expect "login:" { send "root\r" }
expect "#" { send "cd $car_path\r" }
# 解压上传的文件
expect "#" { send "tar -zxvf $file_name\r" }
# 执行完成后保持交互状态
interact
```

执行脚本，`./test.sh my.tar.gz`

## 例子三：混合使用 bash 和 expect

适用于在本地需要计算一些路径或文件名的场景。如以下例子里，需要算出输入文件的 basename，以在跳板机中使用。

```bash
#!/bin/bash

rel_path=$1
file_name=`basename $rel_path`

win_user=mini
win_ip=172.16.1.107
win_pw=win_pw
win_dir=Developer

scp $rel_path $win_user@$win_ip:$win_dir

/usr/bin/expect <<EOF
spawn ssh $win_user@$win_ip
expect {
    "yes/no" { send "yes\r"; exp_continue }
    "password:" { send "$win_pw\r" }
    "#" { send "echo login ok\r" }
}
expect "mini>" { send "adb push $win_dir/$file_name /data\r" }
expect eof
interact
EOF
```

## 参考文章

- [expect - 自动交互脚本](http://xstarcd.github.io/wiki/shell/expect.html)
