# C++ 头文件和库

将实现某部分功能的代码封装成库文件，以方便调用，或者是对代码进行保护加密。

应用场景：有时想将某代码提供给别人用，但是又不想公开源代码，这时可以将代码封装成库文件。在开发中调用其他人员编写的库。

动态库特点：

- 库的代码不会编译进程序里面，所以动态库编译的程序比较小。
- 由动态库编译的程序依赖于系统的环境变量有没有这个库文件，没有则运行不了。

 静态库特点：

- 库的代码会编译进程序里面，所以静态库编译的程序比较大。
- 由静态库编译的程序不用依赖于系统的环境变量，所以环境变量有没有这个库文件，也可以运行。

## 源代码编译成库

- 文件目录结构

```
.
├── my_add.cc
└── CMakeLists.txt
```

- 源文件 my_add.cc

```cpp
int MyAdd(int a, int b) { return a + b; }
```

- CMakeLists.txt

```cmake
# 默认生成静态库
add_library (MyAddStatic my_add.cc)
# 生成动态链接库
add_library (MyAdd SHARED my_add.cc)
```

执行 `cmake . && make` 生成 `libMyAdd.so`（动态库） 和 `libMyAddStatic.a`（静态库）。

> 实际项目中，我们需要把 my_add.cc 的头文件 my_add.h 也要一起编译进来，因为 my_add.h 中可能有全局变量和全局常量的定义，以及一些函数实现。我们下面会进一步解释为什么要这样子做。

## 使用库

如果仅有以上两种格式的文件，我们并不知道里面函数是什么形式的，也不知如何调用他们。为了能让大家使用这些库函数，我们需要编写一个**头文件**，说明这些库里有些什么。对于使用者，只要在程序中声明了该头文件，就可以调用这个库。

头文件只起到了说明的作用（当然有些头文件定义了全局变量、全局常量和函数实现），告诉我们这个库应该怎么调用。所以说我们可以不引头文件，只要知道库的函数是怎么实现的也是可以调用的。

我们把刚刚生成的静态库和动态库拷贝过来，如下目录结构

```
.
├── libMyAdd.so
├── libMyAddStatic.a
└── main.cc
```

- main.cc

```cpp
#include <iostream>

int MyAdd(int a, int b); // 为了调用 MyAdd 函数编译器不报错，需要前置声明

int main() {
  int sum = MyAdd(2, 3);
  std::cout << sum << std::endl; // 输出 5
  return 0;
}

```

使用 gcc 链接 MyAdd 库，`g++ 源文件.cc -L 库路径 -l 库名`

- 链接静态库/动态库。下面是以静态库 `libMyAddStatic.a` 为例子，也可以替换成动态库 `libMyAdd.so`。

```bash
g++ main.cc -L ./ -l MyAddStatic
```

- 运行链接了动态库的程序出现以下报错：

```bash
./a.out: error while loading shared libraries: libMathFunctions.so: cannot open shared object file: No such file or directory
```

解决方式：修改 `LD_LIBRARY_PATH` 环境变量，`export LD_LIBRARY_PATH=$LD_LRBRARY_PATH:动态库的绝对路径`。

## 引用头文件

正常情况下，我们需要为我们的库编写头文件并提供给用户，这样子用户才知道怎么去调用我们的这个库。

- 目录结构

```
.
├── my_add.h
├── libMyAdd.so
├── libMyAddStatic.a
└── main.cc
```

- my_add.h

```cpp
int MyAdd(int a, int b);
```

- main.cc

```cpp
#include <iostream>
#include "my_add.h"

int main() {
  int sum = MyAdd(2, 3);
  std::cout << sum << std::endl; // 输出 5
  return 0;
}

```

## 参考文章

- [C++ 库文件和头文件编写教程](https://mp.weixin.qq.com/s?__biz=MzU0NjgzMDIxMQ==&mid=2247487161&idx=2&sn=400446ef9ac6907499773f0269a667ce&chksm=fb56ec55cc216543e34a8534dcb8b93a7e925060395bb3215dadcf06a8f28aba1756c6178a10&scene=27)

- [Linux C 动态库与静态库的编译与调用](http://t.csdn.cn/UQNbo)