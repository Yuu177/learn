[TOC]

# systemd

Unit

每一个 Unit（服务等） 都有一个配置文件，告诉 Systemd 怎么启动这个 Unit 。

Systemd 默认从目录 /etc/systemd/system/ 读取配置文件。

systemctl reload nginx 重新加载 nginx 配置文件

systemctl restart nginx 重新启动 nginx 服务

## Service

// TODO，service 文件怎么写

## Unit

// TODO

## journal

CentOS 7 以后版本，利用 Systemd 统一管理所有 Unit 的**启动**日志。带来的好处就是，可以只用 journalctl 一个命令，查看所有日志（内核日志和应用日志）。

systemd为我们提供了一个统一中心化的日志系统: journal

其中包含了守护线程 journald 以及我们用来查看日志的工具 journalctl 等等。

journald 任劳任怨，来者不拒地收集来自各个应用和内核的日志信息。

### journal 命令

```bash
# 显示尾部的最新 10 行日志
journalctl -n
# 显示尾部指定行数的日志
journalctl -n 20
# 查看某个 Unit 的日志
journalctl -u nginx.service
# -r reverse 从尾部看日志
journalctl -r
# journalctl 日志太长了会被截断显示不全
journalctl -n 40 -u nginx.service
# 建议使用 vim 打开查看日志
journalctl -n 40 -u nginx.service | vim -
```

参考链接：https://zhuanlan.zhihu.com/p/410995772

### 通过 Systemd Journal 收集日志

systemd journal 收集日志的方式：

- 程序使用库中的 syslog 函数输出的日志使用。
- 任何服务进程输出到 STDOUT/STDERR 的所有内容。

注意，只有以 service 的方式运行程序时，journal 才会捕获 STDOUT/STDERR 输出的内容。像 golang 默认的 print 和 panic 打印的堆栈信息就会被 systemd journal 采集到。

