[TOC]

# Linux 使用手册

## 名词介绍

### PATH

`PATH` 说简单点就是一个字符串变量，当输入命令的时候 LINUX 会去查找 `PATH` 里面记录的路径。所以，path 配置的路径下的文件可以在任何位置执行，并且可以通过 `which` 可执行文件 命令来找到该文件的位置。

- 查看 PATH

`echo $PATH`

- 配置 PATH

> 用 `:` 来分割

```bash
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

### 环境变量

环境变量是进程中一组变量信息，环境变量分为系统环境变量、用户环境变量和进程环境变量。系统有全局的环境变量，在进程创建时，进程继承了系统的全局环境变量、当前登录用户的用户环境变量和父进程的环境变量。进程也可以有自己的环境变量。

### 端口号

- 服务端端口号：我们创建的服务器是需要绑定端口的，这样才能被客户端正确连接。

- 客户端端口号：而客户端在连接后使用的端口号是由操作系统**动态分配**的（当然也可以设置为固定的端口号）

`netstat` 是一个用于显示网络连接、路由表、接口统计信息等的命令。

```bash
netstat -tuln
```

选项说明：

- `-t` 显示 TCP 端口
- `-u` 显示 UDP 端口
- `-l` 仅显示监听状态的端口
- `-n` 以数字形式显示地址和端口号

## 常用命令

### lsof

`lsof -i:80`

查看哪个程序占用了 80 端口号，如果输出结果为空，可能需要加上 sudo

### ps

查看服务器的进程信息

`-e`：表示列出全部的进程
`-f`：显示全部的列（显示全字段）

如查看 kafka 进程的 pid：

```bash
ps -ef | grep kafka 
```

- `netstat` 通过 kafka pid 查看地址和端口号

```bash
netstat -tunlp | grep kafka_pid
```

`-t` 或 `–tcp`：显示 TCP 传输协议的连线状况；

`-u` 或 `–udp`：显示 UDP 传输协议的连线状况；

`-n` 或 `–numeric`：直接使用 ip 地址，而不通过域名服务器；

`-l` 或 `–listening`：显示监控中的服务器的 Socket；

`-p` 或 `–programs`：显示正在使用 Socket 的程序识别码和程序名称；

- `ps aux` 和 `ps -ef`  区别

两者的输出结果差别不大，但展示风格不同。`aux` 是 BSD 风格，`-ef` 是 System V 风格。这是次要的区别，一个影响使用的区别是 **aux 会截断 command 列**，而 `-ef` 不会。当结合 grep 时这种区别会影响到结果。

### ls

`ls -l --block-size=m` 文件大小以 M 显示

### dpkg

[Ubuntu 系统 dpkg 命令使用详解](https://www.jianshu.com/p/2ec0f4b945a2)

安装好后在显示所有应用程序没有显示图标。把相关程序的 Desktop Entry（文件以 `.desktop` 为后缀名） 复制到 `/usr/share/applications` 这个文件目录下。然后点击左下角的显示应用程序就可以看到程序的图标了。

```bash
$ cat com.hmja.notepad.desktop
[Desktop Entry]
Name=Notepad--
Version=1.0
Exec=/opt/apps/com.hmja.notepad/files/Notepad-- %U
Comment=文本编辑器
Icon=/opt/apps/com.hmja.notepad/entries/icons/hicolor/scalable/apps/ndd.svg
Type=Application
Terminal=false
StartupNotify=true
Encoding=UTF-8
Categories=TextEditor;
MimeType=text/plain;text/html;text/x-php;text/x-c;text/x-shellscript;
TryExec=/opt/apps/com.hmja.notepad/files/Notepad--
X-DDE-FileManager-MenuTypes=SingleFile
```

> Desktop Entry 文件是 Linux 桌面系统中用于描述程序启动配置信息的文件。Desktop Entry 文件实现了类似于 Windows 操作系统中快捷方式的功能。

### du

disk usage，显示每个文件和目录的磁盘使用空间

- 命令格式

```bash
du [选项][文件]
```

如果不指定文件（文件夹），则默认为当前目录的占用空间

- 命令参数

`-h` 或 `--human-readable` 以 K，M，G 为单位，提高信息的可读性。

`--max-depth`=<目录层数> 超过指定层数的目录后，予以忽略。

- 例子

```
du -h --max-depth=1
```

### df

disk free，文件系统的磁盘空间占用情况

- 命令参数

`-h` 或 `--human-readable` 以 K，M，G 为单位，提高信息的可读性。

- 例子

```
df -h
```

### nc

netcat 是一个简单的网络工具。比如常用的 telnet 测试 tcp 端口

```bash
telnet 172.16.2.2 23333
Trying 172.16.2.2...
Connected to 172.16.2.2.
Escape character is '^]'.

HTTP/1.1 400 
Server: WebSocket++/0.8.1

Connection closed by foreign host.
```

而 nc 可以支持测试 Linux 的 tcp 和 udp 端口，而且也经常被用于端口扫描。

```bash
nc -v 172.16.2.2 23333
Connection to 172.16.2.2 23333 port [tcp/*] succeeded!
```

## 查看大文件日志

### grep

过滤出关键字相关的日志。

```
grep "关键字" yourlogfile.log > output.log
```

### split

有时候要看完整的日志上下文，可以把大日志文件切分为几个文件，依次打开。

```
split -l 800000 qnxslog59 output59_
```

- `-l`：值为每一输出档的列数大小
- `800000`：行数
- `qnxslog59`：要切割的文件
- `output59_`：输出文件名的前缀

输出：

```
output59_aa  output59_ab  output59_ac
```

## 查看进程内存

在 Linux 系统中，可以通过查看`/proc`文件系统下的特定文件来获取一个进程的内存使用情况。对于每个运行中的进程，系统都会在 `/proc` 目录下创建一个以进程 ID（PID）命名的目录，例如 `/proc/1234`，其中 `1234` 是进程 ID。在这个目录下，有一个名为 `status` 的文件，它包含了进程的状态信息，包括内存使用情况。

```bash
cat /proc/1234/status
```

- `VmPeak`: 进程的峰值虚拟内存使用量，以 kB 为单位。这是进程生命周期中虚拟内存使用量的最大值。
- `VmSize`: 当前进程的虚拟内存使用量，以 kB 为单位。它包括所有的代码、数据和共享库以及页面交换文件（swap）中替换出去的页面。
- `VmHWM`: 进程的高水位标记（High Water Mark）的物理内存使用量，以 kB 为单位。这是进程生命周期中物理内存使用量的最大值。
- `VmRSS`: Resident Set Size，进程当前占用的物理内存（非交换空间）大小，以 kB 为单位。

**使用 top 命令查看的 VSZ（Virtual Memory Size）为虚拟内存，对应 VmSize**

## bash 脚本命令

### 检测变量是否为空

```bash
if [ -z "$build_for" ]; then
  echo "build_for is empty"
else
  echo "build_for is not empty"
fi
```

