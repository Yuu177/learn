[TOC]

# CMake

生成动态库和静态库

```cmake
# 默认生成静态库
ADD_LIBRARY(hello ${LIBHELLO_SRC})
# 上面的代码可以写为
ADD_LIBRARY(hello STATIC ${LIBHELLO_SRC})
# SHARED 表示生成动态库
ADD_LIBRARY(hello SHARED ${LIBHELLO_SRC})
```

## CMake 升级

下载地址：https://cmake.org/download/

![cmake升级](./.CMake.assets/cmake升级.png)

下载最新版本的 [cmake-3.25.2-linux-x86_64.tar.gz](https://github.com/Kitware/CMake/releases/download/v3.25.2/cmake-3.25.2-linux-x86_64.tar.gz)

```bash
# 解压
tar -xzvf cmake-3.25.2-linux-x86_64.tar.gz
# 解压出来的包，将其放在 /opt 目录下，其他目录也可以，主要别以后不小心删了
sudo mv cmake-3.25.2-linux-x86_64 /opt/
# 建立软链接
sudo ln -sf /opt/cmake-3.25.2-linux-x86_64/bin/*  /usr/bin/
# 查看 cmake 版本
cmake --version
```

## make

执行 cmake 命令后，会生成 makefile 文件。make 命令其实是 make all 的省略，即生成所有目标文件。

make 后面跟的是 target，即要编译的目标**（可执行文件，静态库，动态库）**，在 makefile 里面会列出。这个 target 可能也会依赖其他的 target，最终会到依赖的源文件和头文件。

- CMakeLists.txt

```cmake
# 指定生成 MathFunctions 链接库
add_library (MathFunctions ${DIR_LIB_SRCS})
# 指定生成目标
add_executable(demo main.cc)
```

执行 `make`

```bash
⇒  make
[ 50%] Built target MathFunctions
[100%] Built target demo
```

执行 make target 命令生成可执行文件：`make demo`

```bash
⇒  make demo
[100%] Built target demo
```

`make MathFunctions` 生成静态库

```bash
⇒  make MathFunctions
[100%] Built target MathFunctions
```

## CMakeLists.txt 文件参数

### add_compile_definitions

在 CMake 中，`add_compile_definitions` 命令用于向编译器添加宏定义。

语法如下：

```cmake
add_compile_definitions(definition1 [definition2 ...])
```

其中，`definition1`、`definition2` 等表示要添加的宏定义。可以添加多个宏定义，用空格分隔。

例如，要添加一个名为 `ENABLE_DEBUG` 的宏定义，可以使用以下命令：

```cmake
add_compile_definitions(ENABLE_DEBUG)
```

这个命令等同于在源文件中加上以下代码：

```c++
#define ENABLE_DEBUG
```

如果需要定义一个宏的值，可以使用 `=` 运算符：

```cmake
add_compile_definitions(MY_MACRO=42)
```

这个命令等同于在源文件中加上以下代码：

```cpp
#define MY_MACRO 42
```

通过 `add_compile_definitions` 添加的宏定义会作用于整个项目，可以在所有源文件中使用。

### set

`set` 命令则用于设置 CMake 变量的值。

在 CMake 中，变量不需要事先定义就可以直接使用，使用 `set` 命令可以在 CMakeLists.txt 文件中定义新的变量并赋值。如果指定的变量不存在，则 `set` 命令将创建一个新的变量，并将其设置为指定的值。如果变量已经存在，则 `set` 命令将覆盖该变量的当前值。

在 CMake 中，`set(xxx ON)` 的作用是将变量 `xxx` 的值设置为 `ON`。

语法如下：

```cmake
set(xxx ON)
```

其中，`xxx` 表示变量名，`ON` 表示变量的值。

在 CMake 中，变量的值可以是布尔型、字符串型、列表型等。`ON` 是布尔型变量的一种取值，表示变量的值为真。在 CMake 中，布尔型变量还可以取值 `OFF` 表示假。

这个命令的具体作用取决于 `xxx` 变量在 CMakeLists.txt 文件中的上下文。如果 `xxx` 是一个开关变量（类似于编译选项），则将其设置为 `ON` 表示打开该开关，否则可能有其他含义。

例如，如果有一个名为 `ENABLE_DEBUG` 的开关变量，用于控制是否开启调试模式，可以使用以下命令将其设置为打开：

```cmake
set(ENABLE_DEBUG ON)
```

这个命令等同于设置变量 `ENABLE_DEBUG` 的值为 `true`，表示开启调试模式。在 CMakeLists.txt 文件中，可以根据该变量的值来决定是否编译调试模式的代码。

```c++
#if defined(ENABLE_DEBUG)
std::cout << "debug info" << std::endl; 
#endif
```

## cmake 命令参数

### -D

`cmake -D` 命令可以用来定义新的 CMake 变量或者设置已经定义的 CMake 变量的值。如果变量不存在，则会创建一个新的变量，并将其设置为指定的值；如果变量已经存在，则会将其值设置为指定的值。这个选项的语法如下：

```
cmake -D<var>=<value> ...
```

其中，`<var>` 是要设置的变量名称，`<value>` 是变量的值。

在命令行中使用 `-D` 选项时，可以将变量设置为布尔型、字符串型、路径型、列表型等不同类型。如果不指定类型，则默认为字符串型。如果要将变量设置为布尔型，可以使用 `ON`、`OFF`、`TRUE`、`FALSE` 等关键字。

例如，下面的命令将 CMake 变量 `CMAKE_BUILD_TYPE` 的值设置为 `Release`：

```
cmake -DCMAKE_BUILD_TYPE=Release
```

这个命令告诉 CMake 使用 `Release` 构建类型进行构建，这将启用编译器优化并生成优化代码。通常，`CMAKE_BUILD_TYPE` 变量用于控制编译器如何优化代码，并决定编译器是否生成调试符号。

注意，在使用 `-D` 选项设置变量时，变量名和值之间不应该有空格，否则会导致语法错误。

- `-D` 选项可以在生成 CMake 构建系统之前设置变量的值
- `-D` 选项可以覆盖 CMakeLists.txt 文件中定义的变量的

## 参考文章

- [CMake 入门实战](https://www.hahack.com/codes/cmake/)