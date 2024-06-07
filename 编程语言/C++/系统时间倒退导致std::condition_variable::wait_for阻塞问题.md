[TOC]

# 系统时间倒退导致std::condition_variable::wait_for阻塞问题

## 版本信息

QNX 700，gcc version 5.4.0

## 问题现象

系统时间戳倒退，进程阻塞。

## 问题分析

使用 `date` 命令修改系统时间模拟时间戳倒退，稳定复现了此问题。代码中使用 C++ 条件变量 `wait_for`（阻塞当前线程，直到条件变量被唤醒，或直到指定时限时长后），当执行到 `wait_for` 这行代码时候，这个时候如果时间戳倒退，`wait_for` 的超时到期时间会被修改，最终导致进程阻塞（猜测 `wait_for` 使用了系统时间）。

## 源码分析

查看条件变量引用的头文件 `qnx700_1.2.1.c1/target/qnx7/usr/include/c++/v1/condition_variable`，寻找 `wait_for` 函数的最终的调用：

```C++
void __do_timed_wait(unique_lock<mutex>& __lk,
 chrono::time_point<chrono::system_clock, chrono::nanoseconds>) _NOEXCEPT;
```

> `__do_timed_wait` 在 `qnx700_1.2.1.c1/target/qnx7/usr/include/c++/v1/__mutex_base` 这个头文件下。

根据函数签名的入参 `chrono::system_clock` 推测该函数的具体实现用了系统的时间。

> `chrono::steady_clock` 单调时钟，`chrono::system_clock` 系统时钟

## 解决方案

- 方案一：升级 GCC 版本。
- 方案二：使用第三方库。比如使用 `boost::condition_variable` 替换掉相关的代码（boost 在 1.60 修复了该问题）。

## 相关讨论

- [condition_variable does not use monotonic_clock](https://gcc.gnu.org/bugzilla/show_bug.cgi?id=41861)
- https://www.reddit.com/r/cpp/comments/d782lh/after_10_years_the_gccs_stdcondition_variablewait/