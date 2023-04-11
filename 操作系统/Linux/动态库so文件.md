[TOC]

# 动态库 so 文件

## Linux 中命名系统中共享库的规则

```
libname.so.x.y.z
```

- lib：固定代表共享库
- name：共享库名称
- so：固定后缀
- x：主版本号，不同的版本号之间不兼容
- y：次版本号，增量升级，向后兼容
- z：发行版本号，对应次版本的错误修正和性能提升，不影响兼容性

如 `libavcodec.so.58.54.100`

名字叫做 avcodec 的动态库，版本号为 58。

我们一般会有一个没有版本号的软连接来指向该动态库

```
libavcodec.so -> libavcodec.so.58.54.100
```

## 查询动态库需要安装的包

一般我们代码编译的时候可能会出现类似下面的问题

```
error while loading shared libraries: libavcodec.so.58: cannot open shared object file: No such file or directory
```

- 安装 apt-file

```bash
sudo apt install apt-file
```

- 更新 apt-file

```bash
sudo apt-file update
```

- 使用 apt-file 搜索

```bash
➜  ~ sudo apt-file search libavcodec.so  
libavcodec-dev: /usr/lib/x86_64-linux-gnu/libavcodec.so
libavcodec-extra58: /usr/lib/x86_64-linux-gnu/libavcodec.so.58
libavcodec-extra58: /usr/lib/x86_64-linux-gnu/libavcodec.so.58.54.100
libavcodec58: /usr/lib/x86_64-linux-gnu/libavcodec.so.58
libavcodec58: /usr/lib/x86_64-linux-gnu/libavcodec.so.58.54.100
```

- 安装对应包

```bash
sudo apt install libavcodec-dev
```

