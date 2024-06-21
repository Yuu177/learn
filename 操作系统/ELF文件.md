[TOC]

# ELF 文件

ELF 的全称是 Executable and Linking Format，即 “可执行可连接格式”，通俗来说，就是二进制程序。ELF 规定了这二进制程序的组织规范，所有以这规范组织的文件都叫 ELF 文件。是类 Unix 操作系统的二进制文件标准格式。

在 Linux 系统中，一个 ELF 文件主要用来表示 3 种类型的文件：

1. 可执行文件：被操作系统中的加载器从硬盘上读取，载入到内存中去执行。
2. 目标文件（.o）：被链接器读取，用来产生一个可执行文件或者共享库文件。
3. 共享库文件（.so）：在动态链接的时候，由 ld-linux.so 来读取。

ELF 结构

ELF文件主要由四个部分组成 1.ELF头(ELF header），2.程序头表（Program header），3.节（Section），4节头表（Section header table）。

![ELF结构](.ELF文件.assets/ELF结构.jpg)

ELF 文件主要的用途有两个：构建程序，链接成动态库或者是目标文件。运行程序，一般指运行链接好的 `.so` 或者 `.o` 文件。

这个 ELF 文件用作不同用途，文件结构的解析角度就有点不一样，通俗来说，不同用途对需要哪些数据的要求不一样，例如构建（链接）时节表头（Section Table Header）是必须的，但运行时却是可选的，例如运行需要段信息，而链接只有节信息。

- 链接器的主要功能是生成可执行文件。
- 加载器的主要目标是将可执行文件加载到主存中。
- ELF header 描述了文件的总体信息，以及两个 table 的相关信息(偏移地址，表项个数，表项长度)。
- 每一个 table 中，包括很多个表项 Entry，每一个表项都描述了一个 Section/Segment 的具体信息。

## 查看交叉编译文件类型

在嵌入式开发的时候，经常会涉及交叉编译到 armv7、armv8 等平台（通过 `cat /proc/cpuinfo` 或者 `uname -a` 查看 CPU 架构）。通过 `readelf -h` 命令可以查看编译产物的具体信息。

```bash
readelf -h a.out
```

armv7，32 位，系统架构 ARM

```
ELF 头：
  Magic：   7f 45 4c 46 01 01 01 00 00 00 00 00 00 00 00 00 
  类别:                              ELF32
  数据:                              2 补码，小端序 (little endian)
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI 版本:                          0
  类型:                              REL (可重定位文件)
  系统架构:                           ARM
  版本:                              0x1
  入口点地址：                         0x0
  程序头起点：                         0 (bytes into file)
  Start of section headers:          8064 (bytes into file)
  标志：                              0x5000000, Version5 EABI
  Size of this header:               52 (bytes)
  Size of program headers:           0 (bytes)
  Number of program headers:         0
  Size of section headers:           40 (bytes)
  Number of section headers:         13
  Section header string table index: 12
```

armv8，64 位，系统架构 AArch64（ARM）

```
ELF Header:
  Magic:   7f 45 4c 46 02 01 01 00 00 00 00 00 00 00 00 00 
  Class:                             ELF64
  Data:                              2's complement, little endian
  Version:                           1 (current)
  OS/ABI:                            UNIX - System V
  ABI Version:                       0
  Type:                              DYN (Shared object file)
  Machine:                           AArch64
  Version:                           0x1
  Entry point address:               0x2d4f20
  Start of program headers:          64 (bytes into file)
  Start of section headers:          17764176 (bytes into file)
  Flags:                             0x0
  Size of this header:               64 (bytes)
  Size of program headers:           56 (bytes)
  Number of program headers:         7
  Size of section headers:           64 (bytes)
  Number of section headers:         27
  Section header string table index: 26
```

## 参考文章

- https://blog.csdn.net/nirendao/article/details/123883856
- https://www.elecfans.com/d/1908022.html