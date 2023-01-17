# vscode 插件

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

http://t.csdn.cn/jS3uT

