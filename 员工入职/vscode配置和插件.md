[TOC]

# vscode

## 配置文件说明

- setting.json

settings.json 是 VS Code 的配置文件，可以对 VS Code 进行页面风格、代码格式、字体颜色大小等的编辑设置。

setting.json 分【工作区】和【用户】。想象一下你有好几个项目，不同的项目需要不同的配置，所以你可以在每一个项目根路径添加一个 `.vscode/setting.json` 文件来保存对当前工作区的配置。默认打开的设置界面配置的是全局（用户）设置，对每个工作区都会起作用，不过工作区的配置会覆盖全局配置，优先级更高。

- launch.json 和 tasks.json

简单的理解：tasks.json ---> gcc；launch.json ----> gdb。

启动 gdb 调试会话之前需要首先执行 gcc 编译任务。因此，launch.json 有一条配置 preLaunchTask，指向 tasks.json 中的编译任务（label）。

- c_cpp_properties.json

该文件是 C/C++ 这个官方插件的配置。这个文件主要是用于语言引擎的配置，例如：指定 include 路径，智能感知，问题匹配类型等。`Ctrl+Shift+P` 打开 Command Palette，找到并打开 `C/C++: Edit Configurations` 。进行一些配置后，`.vscode` 文件夹下会自动生成此文件。

## vscode 插件

无论是 golang 还是 c++，一直在用 vscode，可以说 vscode 对于我来说已经无法替代了。记录一些自己用着觉得很好用的插件。

### Go

Rich Go language support for Visual Studio Code

golang 插件，非常非常好用。为什么要用 IDE 写代码，这个插件也许是最好的阐述了。这个插件包括了很多 go 的 tools。静态检查，gofmt 等等，能写出更优雅的代码。

### .gitignore Generator

Lets you easily and quickly generate `.gitignore` file for your project using [gitignore.io](https://gitignore.io/) API.

快速生成 .gitignore 文件。

### 翻译(英汉词典)

划词翻译，本地77万词条英汉词典，不依赖任何在线翻译API（速度快），无查询次数限制。

### background

为 vscode 增添一份色彩，Add a lovely background-image to your vscode.

### Fix VSCode Checksums

An extension to to adjust checksums after changes to VSCode core files. Once the checksum changes are applied and VSCode is restarted, all warning about core file modifications will disappear, such as the display of `[Unsupported]` in the title-bar.

因为使用 background 会使得 vscode 显示损坏，强迫症受不了。所以需要用这个插件修复一下。

### Bracket Pair Colorizer

A customizable extension for colorizing matching brackets.

一眼就能看到对应的括号匹配。

ps：已经被 vscode 官方收编，需要勾选一下启用括号对指南

```json
"editor.guides.bracketPairs": "active"
```

### GitLens — Git supercharged

在 vscode 中可视化代码作者身份

### Git History

查看修改文件的历史版本。

[vscode 的 Git History，GitLens — Git supercharged 插件](http://t.csdn.cn/8eWze)

### ~~MySQL~~

> ~~插件商店中有好几个同名的 MySQL 插件，认准链接：https://marketplace.visualstudio.com/items?itemName=formulahendry.vscode-mysql~~ 

~~适合轻度使用 mysql 的用户，不用再每次都要在终端中用命令行登录 mysql，相当方便。~~

~~目前无法 rename connections，可以通过改 hosts 的方法去解决：https://github.com/formulahendry/vscode-mysql/issues/16~~

因为查询出来的数据太大的话精度会丢失，不再推荐该插件

### Project Manager

Easily switch between projects.

vscode 每次切换目录都很麻烦，有了这个插件都非常快速的切换不同的项目文件。

右键项目配置标签，可以让项目按照标签分类。

### Todo Tree

Show TODO, FIXME, etc. comment tags in a tree view

### Shades of Purple

好看的一个主题

### Bookmarks

经常在文档几个不同位置跳来跳去的，这个书签插件很实用，右键菜单里直接设置/取消书签，快捷键在不同的书签位置跳转，左边还有当前书签列表，双击立马跳转。

### vscode-icons

丰富 vscode 文件 icons 显示

### Carbon

Carbon 能够轻松地将你的源码生成漂亮的图片并分享

https://github.com/carbon-app/carbon/blob/main/docs/README.cn.zh.md

### clang-format

format C++ 代码格式为 Google Style

```bash
sudo apt install clang-format
```

### cpplint

cpplint 是 Google 开发的一个 C++ 代码风格检查工具，遵循 google code style。

filter 参数的用法，就是以 `+` 或者 `-` 开头接着写规则名，就表示启用或者屏蔽这些规则。

```json
{
    "cpplint.filters": [
        "-build/include_subdir", // 我们将有自己的头文件路径规范
        "-whitespace/line_length", // 可选项，按个人喜好来
        "-runtime/references", // 不做要求，该要求原为：函数参数列表必须 (const &) 形式 或者 (* 指针) 形式。
        "-build/c++11", // cpp11 不支持的头文件的提示（老是提示 <thread> 等头文件有问题，故屏蔽）
        "-build/header_guard", // 找不到配置 <PROJECT>_<PATH>_<FILE>_ 中 <PROJECT>_ 的方法，故屏蔽
        "-legal/copyright",
    ],
}
```

### clangd

vscode 官方的 cpptools 在大型 C++ 项目中函数跳转很慢，所以改使用 clangd 代替。

安装 clangd 插件：https://zhuanlan.zhihu.com/p/364518020。

而 clangd 是基于 `compile_commands.json` 文件来完成对项目的解析，并支持代码补全和跳转。

该文件有三种生成方式：

- 使用 cmake 生成。正常执行 cmake 命令，我们会生成一个 build 目录。使用 cmake 生成 compile_commands.json，需要在运行 cmake 时添加参数 `-DCMAKE_EXPORT_COMPILE_COMMANDS=True` 或者在 CMakeLists.txt 中添加 `set(CMAKE_EXPORT_COMPILE_COMMANDS True)`。这样子我们在 build 目录下就会看到一个 compelie_commands.json 文件了。

- 如果是基于 make 方式来编译，那么可以先安装 `pip install compiledb`，之后在当前目录下运行

  - `compiledb -n make -C build`
  - `compiledb make -C build`

  这两个命令中的其中一个来生成 compile_commands.json 文件，其中前者不会执行真正的 make 编译命令。

- 如果是基于其他方式，可以使用 https://github.com/rizsotto/Bear 项目中的方式来生成对应的 compile_commands.json 文件。

在 setting.json 下配置：

生成 compile_commands.json 文件后，我们只需要配置 `--compile-commands-dir` 来指定 compile_commands.json 所在的目录即可。

```json
{
    // clangd 位置，使用 vscode 插件商店下载会自动配置
    "clangd.path": "/home/tanpanyu/.config/Code/User/globalStorage/llvm-vs-code-extensions.vscode-clangd/install/15.0.6/clangd_15.0.6/bin/clangd",
    "clangd.arguments": [
        // 在后台自动分析文件（基于 complie_commands)
        "--background-index",
        // 标记 compelie_commands.json 文件的目录位置
        "--compile-commands-dir=${workspaceFolder}/build",
        // 同时开启的任务数量
        "-j=12",
        // clang-tidy 功能
        "--clang-tidy",
        "--clang-tidy-checks=performance-*,bugprone-*",
        // 全局补全（会自动补充头文件）
        "--all-scopes-completion",
        // 更详细的补全内容
        "--completion-style=detailed",
        // 补充头文件的形式
        "--header-insertion=iwyu",
        // pch 优化的位置
        "--pch-storage=disk",
    ],
    // 使用 clangd 会和 vscode 默认推荐的 C/C++ 插件的 IntelliSense 冲突，需要禁止。
    "C_Cpp.intelliSenseEngine": "disabled",
}
```

- 交叉编译找不到标准库头文件

使用交叉编译，在 x86 上编译 ARM 程序。然而 clangd 插件会提示找不到标准库的头文件，因为交叉编译使用的是交叉编译器，它有自己的标准库头文件目录。

解决：在 CMakeLists.txt 下添加

```cmake
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
if(CMAKE_EXPORT_COMPILE_COMMANDS)
    set(CMAKE_CXX_STANDARD_INCLUDE_DIRECTORIES ${CMAKE_CXX_IMPLICIT_INCLUDE_DIRECTORIES})
endif()
```

参考文章：[交叉编译，clang-tidy 找不到交叉编译的标准库头文件](http://t.csdn.cn/keJ4i)

> clang-tidy 是一个基于 clang 的 C++ 静态分析工具，主要用来检测代码中的常见错误

**clangd 配置**

there are two ways to change the [clangd config options](https://clangd.llvm.org/config.html):

- creating a `~/.config/clangd/config.yaml` text file, which will affect all projects
- creating a `.clangd` in a specific project directory

可以通过 vscode 命令来创建

![clangd配置文件](./.vscode配置和插件.assets/clangd配置文件.png)

编辑配置文件文件，如：屏蔽特定告警等。

```yaml
Diagnostics:
  Suppress: 'builtin_definition' # 屏蔽该 builtin_definition 告警
CompileFlags:
  Add: -ferror-limit=0 # Too many errors emitted, stopping now  [clang: fatal_too_many_errors] 消除这个告警
```

参考文章

- [Clangd config](https://ahmadsamir.github.io/posts/12-clangd-config-tweaks.html)
- [Too many errors emitted, stopping now [clang: fatal_too_many_errors]](https://github.com/clangd/coc-clangd/issues/255)



