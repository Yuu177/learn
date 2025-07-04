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

- make 命令

`make distclean` 用于清理源代码目录中的所有生成文件和配置文件，使其恢复到最初的状态。这个命令通常用于确保在重新配置和编译之前，所有旧的生成文件和配置文件都已被删除，以避免潜在的冲突和错误。

`make distclean` 会执行以下操作：

1. 删除所有编译生成的目标文件（例如 `.o` 文件）。
2. 删除所有生成的可执行文件和库文件。
3. 删除所有由 `configure` 脚本生成的配置文件（例如 `config.h`）。
4. 删除其他临时文件和目录。

`make distclean` 通常比 `make clean` 更彻底，因为 `make clean` 可能只删除编译生成的目标文件和可执行文件，而不会删除配置文件。

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

### list

[cmake 命令之 list 介绍](https://www.jianshu.com/p/89fb01752d6f)

### INCLUDE_DIRECTORIES & ADD_SUBDIRECTORY

`INCLUDE_DIRECTORIES` 和 `ADD_SUBDIRECTORY` 是 CMake 中两个不同的指令，用于不同的目的。

`INCLUDE_DIRECTORIES` 指令用于向 CMake 添加一个或多个包含目录，以便编译器可以在编译时找到头文件。例如，假设你有一个名为 `my_project` 的 C++ 项目，并且它的头文件在 `include` 目录中。你可以使用以下代码将该目录添加到 CMake 构建中：

```cmake
INCLUDE_DIRECTORIES(include)
```

`ADD_SUBDIRECTORY` 指令用于向 CMake 添加一个子目录，该子目录包含要构建的另一个 CMake 项目。该指令用于将一个项目的构建过程分成多个子目录，以便可以更轻松地管理和构建代码库。例如，假设你的 `my_project` 依赖于 `my_library` 库。你可以使用以下代码将 `my_library` 添加到 `my_project` 的构建过程中：

```cmake
ADD_SUBDIRECTORY(my_library)
```

这将使 CMake 进入 `my_library` 目录，并使用 `CMakeLists.txt` 文件中的指令构建该库。然后，它会回到 `my_project` 目录，继续构建 `my_project`。

因此，`INCLUDE_DIRECTORIES` 用于指定要包含的头文件的路径，而 `ADD_SUBDIRECTORY` 用于将另一个 CMake 项目添加到当前项目的构建中。

- 那为什么有些项目没有使用到 `INCLUDE_DIRECTORIES` 也能编译通过？

```
./Demo
    |
    +--- main.cc
    |
    +--- math/
          |
          +--- MathFunctions.cc
          |
          +--- MathFunctions.h
```

项目根目录 Demo 和 math 目录里各编写一个 CMakeLists.txt 文件

根目录中的 CMakeLists.txt ：

```cmake
# CMake 最低版本号要求
cmake_minimum_required (VERSION 2.8)
# 项目信息
project (Demo)
# 查找当前目录下的所有源文件
# 并将名称保存到 DIR_SRCS 变量
aux_source_directory(. DIR_SRCS)
# 添加 math 子目录
add_subdirectory(math)
# 指定生成目标
add_executable(Demo main.cc)
# 添加链接库
target_link_libraries(Demo MathFunctions)
```

子目录 math 中的 CMakeLists.txt：

```cmake
# 查找当前目录下的所有源文件
# 并将名称保存到 DIR_LIB_SRCS 变量
aux_source_directory(. DIR_LIB_SRCS)
# 生成链接库
add_library (MathFunctions ${DIR_LIB_SRCS})
```

在这种情况下，`MathFunctions.h` 头文件位于 `math` 目录中，并且使用了 `aux_source_directory` 命令来将 `math` 目录下的所有源文件添加到变量 `DIR_LIB_SRCS` 中，然后使用 `add_library` 命令来创建名为 `MathFunctions` 的库。

由于 `MathFunctions.h` 头文件包含在 `math` 目录下的源文件中，因此 `DIR_LIB_SRCS` 变量会包含 `MathFunctions.h` 头文件的路径，这样编译器就可以在编译时找到头文件。但是，如果头文件位于其他目录中，就需要使用 `INCLUDE_DIRECTORIES` 命令来将包含路径添加到 CMake 构建中。

在上述情况中，因为头文件位于 `math` 目录中，所以不需要使用 `INCLUDE_DIRECTORIES` 命令。但是，如果头文件位于其他目录中，应该使用类似于以下的命令将包含路径添加到 CMake 构建中：

```cmake
INCLUDE_DIRECTORIES(include)
```

其中，`include` 是头文件所在的目录的路径。使用此命令将 `include` 目录添加到 CMake 构建中，编译器就可以找到头文件，并在编译时包含它们。

### link_directories 和 include_directories 区别

`link_directories` 和 `include_directories` 分别用于指定编译器在编译和链接过程中查找库和头文件的位置。

`include_directories` 用于指定头文件的搜索路径，即编译器在编译过程中应该搜索的路径。可以使用绝对路径或相对路径，也可以使用 CMake 变量，例如 `${PROJECT_SOURCE_DIR}/include`。

例如，下面的代码将在编译过程中添加一个名为 `my_include_dir` 的搜索路径：

```cmake
include_directories(${my_include_dir})
```

`link_directories` 用于指定库文件的搜索路径，即编译器在链接过程中应该搜索的路径。与 `include_directories` 类似，可以使用绝对路径或相对路径，也可以使用 CMake 变量，例如 `${PROJECT_BINARY_DIR}/lib`。

例如，下面的代码将在链接过程中添加一个名为 `my_library_dir` 的搜索路径：

```cmake
link_directories(${my_library_dir})
```

需要注意的是，虽然 `link_directories` 可以用于指定链接时查找库文件的路径，但是更好的做法是使用 `target_link_directories` 命令，该命令可以将库路径与特定的目标进行关联。

例如，下面的代码将在链接 `my_target` 时添加一个名为 `my_library_dir` 的搜索路径：

```cmake
target_link_directories(my_target PRIVATE ${my_library_dir})
```

总的来说，`include_directories` 用于指定头文件的搜索路径，`link_directories` 用于指定库文件的搜索路径，而 `target_link_directories` 用于将库路径与特定的目标进行关联。

### 设置输出的文件目录

```cmake
SET(EXECUTABLE_OUTPUT_PATH "${PROJECT_BINARY_DIR}/bin")
SET(LIBRARY_OUTPUT_PATH "${PROJECT_BINARY_DIR}/lib")
```

### find_package

> 用于查找并配置整个软件包。一个软件包可能包含多个库、头文件、配置文件等。通常依赖于 CMake 提供的模块或包提供的配置文件来完成查找和配置。

#### 查找模式

- Module 模式：这个模式下 CMake 会去 `CMAKE_MODULE_PATH` 找 `Find<PackageName>.cmake`，这种一般是 CMake 或库的用户提供的。
- Config 模式：这个模式下 CMake 会找 `<package_name>-config[-version].cmake` 或 `<PackageName>Config[Version].cmake` ，[查找的目录更加细致 ](https://cmake.org/cmake/help/latest/command/find_package.html#search-procedure)，这种一般由库提供。

没有指定模式的情况下 CMake 会先使用 Module 模式，失败后 fallback 到 Config 模式。如果需要指定，可以：

```bash
find_package(OpenCV REQUIRED)     # 没有指定模式
find_package(package_name MODULE) # 仅使用 Module 模式，不 fallback 到 Config 模式
find_package(package_name CONFIG) # 直接使用 Config 模式
```

**`[REQUIRED]`: 如果指定，表示这个包是必须的，如果找不到会导致 CMake 出错**。

#### 查找路径

```bash
CMake Error at CMakeLists.txt:13 (find_package):
  By not providing "FindOpenCV.cmake" in CMAKE_MODULE_PATH this project has
  asked CMake to find a package configuration file provided by "OpenCV", but
  CMake did not find one.

  Could not find a package configuration file provided by "OpenCV" with any
  of the following names:

    OpenCVConfig.cmake
    opencv-config.cmake

  Add the installation prefix of "OpenCV" to CMAKE_PREFIX_PATH or set
  "OpenCV_DIR" to a directory containing one of the above files.  If "OpenCV"
  provides a separate development package or SDK, be sure it has been
  installed.


-- Configuring incomplete, errors occurred!
```

假设 `OpenCVConfig.cmake` 在 `/tmp/opencv` 这个目录下，通过设置 `<PackageName>_DIR` 或者 `CMAKE_PREFIX_PATH` 让 cmake 找到包配置文件：

```cmake
set(OpenCV_DIR /tmp/opencv)
# or
set(CMAKE_PREFIX_PATH /tmp/opencv)
```

#### 参考

https://blog.csdn.net/zhanghm1995/article/details/105466372

### target_link_libraries

如果库 A 依赖于库 B，那么在链接时，库 B 必须在库 A 之后。例如，如果 `libA` 使用了 `libB` 中的符号，链接器需要先找到 `libB` 中的这些符号，然后再处理 `libA`。

### find_library

`find_library` 用于查找单个库文件（通常是 `.lib`, `.a`, 或 `.so` 文件）。它的作用是找到指定的库文件，并将其路径存储在一个变量中。

1. `EXAMPLE_LIB`：这是你要存储找到的库路径的变量名称。如果找到了库，库的完整路径将被存储在这个变量中。
2. `example`：这是你要查找的库的名称。CMake 将在指定的路径中查找名为 `libexample.so`
3. `PATHS` ：指定的查找路径。CMake 将在这个目录中查找库文件。

```cmake
find_library(EXAMPLE_LIB example PATHS /first/path /second/path /third/path)
if(NOT EXAMPLE_LIB)
    message(FATAL_ERROR "Could not find the example library")
endif()
target_link_libraries(my_executable ${EXAMPLE_LIB})
```

## cmake 命令参数和变量

### -D

`cmake -D` 命令可以用来定义新的 CMake 变量或者设置已经定义的 CMake 变量的值（相当于代码中的全局变量）。如果变量不存在，则会创建一个新的变量，并将其设置为指定的值；如果变量已经存在，则会将其值设置为指定的值。这个选项的语法如下：

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

## 编译选项 -fPIC

```cmake
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fPIE -pie")
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -fPIE -pie")
```

- `CMAKE_CXX_FLAGS` 是用来添加 C++ 编译器标志的变量。
- `CMAKE_C_FLAGS` 是用来添加 C 编译器标志的变量。

`-fPIC/-fpic` ：编译选项，用于生成位置无关的代码（Position-Independent-Code），用于生成动态库。

`-fPIE`：与 `-fPIC` 类似，差别就是生成的 `.o` 文件不能用来链接生成动态库，只能用来生成可执行文件。

`-pie`：是一个链接选项，它要求链接器使用的所有 `.o` 文件编译是必须使用 `-fPIC` 或者 `-fPIE`。

### 地址无关代码

《程序员的自我修养》中「7.3 地址无关代码」有详细介绍

> 补充说明「指令部分」：在动态链接的上下文中，「指令部分」通常指的是程序中的代码段（也称为文本段或 `.text` 段），这是包含了程序执行的机器指令的内存区域。当一个可执行文件或共享库被装载到内存时，它的代码段会被放置到进程的地址空间中。

## cmake-presets

One problem that CMake users often face is sharing settings with other people for common ways to configure a project

https://cmake.org/cmake/help/latest/manual/cmake-presets.7.html

## 参考文章

- [CMake 入门实战](https://www.hahack.com/codes/cmake/)