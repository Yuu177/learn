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

## 待补充

## 参考文章

- https://blog.csdn.net/nirendao/article/details/123883856
- https://www.elecfans.com/d/1908022.html