[TOC]

# QNX

QNX 是一种商用的类 Unix 实时操作系统，遵从 POSⅨ 规范，目标市场主要是嵌入式系统。

## 常用命令

### hogs 查看进程 CPU 和内存占用

http://www.qnx.com/developers/docs/6.5.0/index.jsp?topic=%2Fcom.qnx.doc.neutrino_utilities%2Fh%2Fhogs.html

显示所有进程对系统 CPU 和内存的占用情况。`hogs <pid>` 只查看对应的 pid，或者 `hogs | grep QDms` 查看特定的进程信息。

```shell
# hogs | grep QDms
      PID           NAME  MSEC PIDS  SYS       MEMORY
279535725    a.out     1   0%   0%   3772k   0%
279539827    a.out  1465   6%  48%  41312k   4%
```

- `PID`：进程 ID
- `NAME`：进程名
- `MSEC`：自上次迭代以来此进程运行的毫秒数
- `PIDS`：进程在此迭代中，占所有进程运行时间的百分比
- `SYS`：进程在此迭代中，CPU 的百分比（**分别在不同的 CPU 核心上的占用率的和**）
- `MEMORY`：内存使用情况

>The `SYSTEM` column is incorrect on multicore systems; the numbers in this column will add up to (roughly) the number of processors times 100%. Use the [`top`](http://www.qnx.com/developers/docs/6.5.0/topic/com.qnx.doc.neutrino_utilities/t/top.html) utility instead.

这个 CPU 占用率（SYS）在多核系统不是很准确。比如我有两个核心，DMS 进程在 Processor1 上占用 `10%`，在 Processor2 占用 `10%`，那么 `SYS` 会显示 `20%（20% / 200%）`。但是小核占 `10%` 和大核占 `10%` 是不一样的，所以简单的相加会有一定的误差。

#### 查看 CPU 核心数

通过 `pidin info` 查看 CPU 核心数。

```bash
# pidin info 
CPU:AARCH64 Release:7.0.X  FreeMem:976MB/15GB BootTime:Jan 01 08:00:01 CST 1970
 Actual resident free memory:917Mb
Processes: 118, Threads: 1120
Processor1: 1373601886 Kryo v4 Silver 1612MHz FPU 
Processor2: 1373601886 Kryo v4 Silver 1612MHz FPU 
Processor3: 1373601886 Kryo v4 Silver 1612MHz FPU 
Processor4: 1373601886 Kryo v4 Silver 1612MHz FPU 
Processor5: 1373601870 Kryo v4 Gold 2131MHz FPU 
Processor6: 1373601870 Kryo v4 Gold 2131MHz FPU 
Processor7: 1373601870 Kryo v4 Gold 2131MHz FPU 
Processor8: 1373601870 Kryo v4 Gold Plus 2419MHz FPU 
# 
```

通过上面打印，可以确定该系统的 CPU 是有 8 个核

```bash
# hogs | grep QDms
      PID           NAME  MSEC PIDS  SYS       MEMORY
279539827    a.out  1465   6%  48%  41312k   4%
```

所以这里是 `48% / 800%`

### top 查看 CPU 占用率

```bash
# top -h
top: illegal option -- h
top - display system usage (UNIX)

top  top [-i <number>] [-d] [-n <node>]
Options:
 -d         dumb terminal
 -b         batch mode for background operation
 -n <node>  remote node
 -p <pri>   run at priority
 -i <iter>  # of iterations
 -z <num>   number of threads to display
 -D <delay> delay in seconds
 -t         display thread names
```

`top` 显示 CPU 占用率最高的 10 个线程，可以结合 grep 命令查询。`top -z 40` 显示 CPU 占用率最高的前 40 个线程。

> 一个进程可能有多个线程，把该进程所有的线程的 CPU 占用率加起来就是该进程的 CPU 占用率

```bash
top -b -z 40 -t 1 -D 2
```

### showmem 查看进程内存占用

查看所有进程的内存占用情况。`showmem -p <pid>` 只查看对应的进程。

### pidin

查看进程等信息，类似 `ps` 命令

```bash
pidin | grep test
```

