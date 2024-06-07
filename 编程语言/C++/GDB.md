[TOC]

# GDB

## 调试目标程序

1. 编译代码时添加调试信息：在编译代码时，使用 `-g` 选项来添加调试信息。这将确保生成的可执行文件包含符号表和调试信息，以便 GDB 可以准确地对应源代码和二进制代码。

2. 启动 GDB：在终端中执行以下命令，启动 GDB 并加载要调试的可执行文件：

   ```
   gdb <executable>
   ```

   其中 `<executable>` 是要调试的可执行文件的路径。

3. 设置断点：使用 `break` 命令在代码中设置断点，以在特定位置暂停程序的执行。例如，使用以下命令在函数的开头设置断点：

   ```
   break function_name
   ```

   或者，使用以下命令在指定的源代码行设置断点：

   ```
   break file_name:line_number
   ```

4. 启动程序：使用 `run` 命令启动程序并开始调试。如果需要传递命令行参数，可以在 `run` 命令后面添加参数。例如：

   ```
   run arg1 arg2
   ```

5. 执行调试命令：一旦程序开始执行，可以使用以下常见的 GDB 调试命令来操作程序的执行和状态：

   - `next`（简写为 `n`）：单步执行，逐行执行代码。
   - `step`（简写为 `s`）：单步执行，进入函数内部。
   - `continue`（简写为 `c`）：继续执行程序，直到下一个断点或程序结束。
   - `print`（简写为 `p`）：打印变量的值。
   - `backtrace`（简写为 `bt`）：打印当前的函数调用栈。
   - `frame`（简写 `f`）：切换到对应的帧。
   - `watch`：设置监视点，当变量的值发生变化时暂停程序的执行。
   - `info`：获取关于程序状态的信息，如当前的堆栈帧、变量列表等。比如 `info threads`，`info frame`。
   - `list`（简写 `l`）：显示当前行的源码。
   - `thread`：切换到对应的线程。
   - `quit`（简写为 `q`）：退出 GDB。

## 解析 core 文件

Coredump 是指程序运行时发生错误，导致程序崩溃时，操作系统将当前进程的内存状态信息保存到一个称为 core 文件的特殊文件中。通过分析 core 文件，可以找到导致程序崩溃的原因。

```
gdb <executable> <core>
```

然后使用命令 `bt` 或 `where` 来查看程序崩溃时的函数调用栈，以确定崩溃位置。

因为我们的程序有时候会依赖动态库，如果不引入这些动态库，就可能无法看到堆栈信息。

```
(gdb) bt
#0  0x00000033f50e6510 in ?? ()
#1  0x00000033f50cbd50 in ?? ()
Backtrace stopped: previous frame identical to this frame (corrupt stack?)
```

通过 `info sharedlibrary` 查看需要引入哪些动态库

```
(gdb) info sharedlibrary
From                To                  Syms Read   Shared Object Library
0x00000033f5096000  0x00000033f5096000  No          libc.so.4
```

通过 `set solib-search-path` 设置依赖的动态库目录路径，把动态库加载进来。如果要设置多个路径，可以用过 `:` 分隔。

### QNX 解析 core

QNX 系统解析 core 文件需要用到**程序发生 core dumped 时候使用的 so 库**，不然就会出现无法把动态库加载进来的情况。

```bash
warning: Shared object "/home/jinx/code/toolchain/qnx700/target/qnx7/aarch64le/lib/libc.so.4" could not be validated and will be ignored.

warning: Could not load shared library symbols for 13 libraries, e.g. /usr/lib/ldqnx-64.so.2.
```

### 交叉编译工具链解析 core

cd 到工具链目录下，切换到 bash，然后执行工具链设置环境变量的脚本。这样子我们使用 gdb 加载 core 文件的时候就不需要手动执行 `set solib-search-path` 来指定要加载的动态库路径了。

下面以 QNX700 举例如何设置工具链的环境变量：

```bash
➜  ls
bsp  buildinfo  custom  host  license  qnxsdp-env.sh  Readme.txt  target
➜  bash
$ source qnxsdp-env.sh
QNX_HOST=/home/jinx/code/toolchain/qnx700_1.2.1.c1/host/linux/x86_64
QNX_TARGET=/home/jinx/code/toolchain/qnx700_1.2.1.c1/target/qnx7
MAKEFLAGS=-I/home/jinx/code/toolchain/qnx700_1.2.1.c1/target/qnx7/usr/include
$
```

## Linux 设置生成 core 文件

### 内核参数

内核参数是可在系统运行时调整的可调整值。不需要重启或重新编译内核就可以使更改生效。

`sysctl` 命令用于运行时配置内核参数，这些参数位于 `/proc/sys` 目录下（先决条件：根权限）。

常用参数的意义：

- -`w`：临时改变某个指定参数的值，如 `sysctl -w net.ipv4.ip_forward=1`
- `-a`：显示所有的系统参数
- `-p`：从指定的文件加载系统参数，如不指定即从 `/etc/sysctl.conf` 中加载

如果仅仅是想临时改变某个系统参数的值，可以用两种方法来实现，例如想启用 IP 路由转发功能：

- `echo 1 > /proc/sys/net/ipvsysctl4/ip_forward`
- `sysctl -w net.ipv4.ip_forward=1`

如果想永久保留配置，可以修改 `/etc/sysctl.conf` 文件，将 `net.ipv4.ip_forward=0` 改为 `net.ipv4.ip_forward=1`

### 核心转储限制

1. 确认核心转储（core dump）的限制：首先，确保核心转储的限制处于适当的状态。在终端中执行以下命令，检查 `ulimit` 设置：

```
ulimit -a
```

确认 `core file size` 的值不是 `0`，而是一个非零值（以 KB 为单位），表示核心转储文件的最大大小限制。如果该值为 `0`，则需要更改核心转储限制。

#### 临时修改

> 只在当前终端窗口生效

```bash
ulimit -c unlimited
```

#### 永久修改

设置核心转储限制：打开 `/etc/security/limits.conf` 文件以编辑，执行以下命令：

```bash
sudo vi /etc/security/limits.conf
```

在文件末尾添加以下行：

```
*    soft    core    unlimited
*    hard    core    unlimited
```

这将设置所有用户（`*`）的核心转储限制为无限制。

### 指定生成 core 文件的路径和文件名模式

在 Linux 上，生成的 core 文件的路径由 `core_pattern` 内核参数指定。可以通过以下方式来确定生成 core 文件的路径：

1. 查看 `core_pattern` 的值：

```bash
cat /proc/sys/kernel/core_pattern
```

该命令将显示当前的 `core_pattern` 值，它指定了生成 core 文件的路径和文件名模式。如 `/cores/core.%e.%p.%t`，表示 core 文件将被保存在 `/cores` 目录下，并使用 `%e`（可执行文件名）、`%p`（进程 ID）和 `%t`（时间戳）作为文件名的一部分。

#### 临时修改

```bash
sudo sysctl -w kernel.core_pattern=/tmp/core.%e.%p.%t
```

将 `/tmp` 替换为想要保存 core 文件的路径。修改 `core_pattern` 后，重新运行程序并让其生成 core 文件。core 文件将保存在指定的路径下。

> 不加目录的话，指定的 core 文件路径就是当前发生 coredump 的程序目录。

#### 永久修改

更新系统配置：打开 `/etc/sysctl.conf` 文件以编辑，执行以下命令：

```
sudo vi /etc/sysctl.conf
```

在文件末尾添加以下行：

```
kernel.core_pattern = core
```

这将指定生成的核心转储文件的名称为 "core"。

重新加载 sysctl 配置：执行以下命令以重新加载 sysctl 配置：

```
sudo sysctl -p
```

现在，在程序崩溃时，应该会生成相应的 core 文件。

请注意，生成 core 文件的路径可能受到系统的设置和权限限制。确保对所选路径具有适当的写入权限，并考虑磁盘空间的大小和可用性，以避免 core 文件填满磁盘。

## 远程调试

TODO

QNX
