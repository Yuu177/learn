# pstack

pstack 的工作原理其实就是一个 shell 脚本，然后在脚本里面调用 gdb 来实现对应用进程各个线程堆栈的打印。

安装 pstack 后，使用 pstack pid 报错

>pid 查询， ps aux | grep xx

```bash
pstack: Input/output error
failed to read target.
```

使用 `which pstack` 查看 pstack 安装在哪里（/usr/bin/pstack）

`cat /usr/bin/pstack` 发现是一堆乱码

之前拷贝了网上的脚本，pstack 后发现打印的是 `?? ()`

```bash
#2  0x00007ffcd186d57e in ?? ()
#3  0x000055de22d78080 in ?? ()
#4  0x0000000000000000 in ?? ()
```

后面从另一台 centos 拷贝过来 pstack 脚本可以正常使用。

```bash
#!/bin/sh

if test $# -ne 1; then
    echo "Usage: `basename $0 .sh` <process-id>" 1>&2
    exit 1
fi

if test ! -r /proc/$1; then
    echo "Process $1 not found." 1>&2
    exit 1
fi

# GDB doesn't allow "thread apply all bt" when the process isn't
# threaded; need to peek at the process to determine if that or the
# simpler "bt" should be used.

backtrace="bt"
if test -d /proc/$1/task ; then
    # Newer kernel; has a task/ directory.
    if test `/bin/ls /proc/$1/task | /usr/bin/wc -l` -gt 1 2>/dev/null ; then
        backtrace="thread apply all bt"
    fi
elif test -f /proc/$1/maps ; then
    # Older kernel; go by it loading libpthread.
    if /bin/grep -e libpthread /proc/$1/maps > /dev/null 2>&1 ; then
        backtrace="thread apply all bt"
    fi
fi

GDB=${GDB:-gdb}

# Run GDB, strip out unwanted noise.
# --readnever is no longer used since .gdb_index is now in use.
$GDB --quiet -nx $GDBARGS /proc/$1/exe $1 <<EOF 2>&1 |
set width 0
set height 0
set pagination no
$backtrace
EOF
/bin/sed -n \
    -e 's/^\((gdb) \)*//' \
    -e '/^#/p' \
    -e '/^Thread/p'
```