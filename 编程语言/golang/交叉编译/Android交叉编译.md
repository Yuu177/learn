[TOC]

# Go 交叉编译 Android 可执行程序

## 一般交叉编译

一般在 Linux 平台上进行开发，但是最近有个需求需要在安卓上启动一个 http 服务器。

GO 支持交叉编译，编译成安卓程序用以下命令即可：

```bash
GOOS=android GOARCH=arm64 go build
```

用上面的命令生成的可执行文件在安卓上能够正常运行，但是编译出来的程序无法解析 `DNS`（在 `/etc/hosts` 添加对应的域名映射后程序能够正常对 `DNS` 解析），报错如下：

```bash
Get "https://www.baidu.com": dial tcp: lookup www.baidu.com on [::1]:53: read udp [::1]:42913->[::1]:53: read: connection refused
```

在 GitHub 搜索后，发现了相同的问题：[[Bug]: Problems with DNS (?) when trying to send email with shoutrrr](https://github.com/termux/termux-app/issues/3156#issuecomment-1396136966)

问题原因：Go apps need to be compiled with **cgo** and with android set as the goos because termux doesn't have a resolve.conf.

## 交叉编译 CGO

### 下载 NDK

>Android NDK 是一个工具集，可让您使用 C 和 C++ 等语言以原生代码实现应用的各个部分。对于特定类型的应用，这可以帮助您重复使用以这些语言编写的代码库。

交叉编译 Android 版本的动态库不仅需要指定 `GCC`，还要指定 `NDK_TOOLCHAIN`，所以要先下载对应的 [NDK](https://github.com/android/ndk/wiki/Unsupported-Downloads)。我这里下载的是 `r20` 版本。

### 编译 toolchain

> 下载解压后的文件夹名 `android-ndk-r21e`

我们使用 `ndk` 自带的 `./android-ndk-r21e/build/tools/make-standalone-toolchain.sh` 脚本，编译特定的 `toolchain`，使用命令如下：

```bash
./make-standalone-toolchain.sh --toolchain=aarch64-linux-android-4.9 --platform=android-30 --install-dir=~/toolchains/ndk_toolchain
```

**脚本参数说明：**

- `toolchain`：表示对应 Android 的 `ARCH`。`arm32` 使用 `arm-linux-androideabi-4.9`，`arm64` 使用 `aarch64-linux-android-4.9`。对应的是 android-ndk-r21e 文件夹下的 toochains。如果要交叉编译 x86 架构，注意生成 --toolchain=x86_64-4.9

```bash
android-ndk-r21e/toolchains
├── aarch64-linux-android-4.9
├── arm-linux-androideabi-4.9
├── llvm
├── renderscript
├── x86-4.9
└── x86_64-4.9
```

- `platform`：表示对应 Android API 的版本。`android-30` 对应 `Android 11` 系统。更多请参考：[Android 平台版本所支持的 API 级别](https://developer.android.google.cn/guide/topics/manifest/uses-sdk-element?hl=zh-cn)
- `install-dir`：表示编译的目标 `toolchain` 存放位置，后面交叉编译 Go 代码会用到。

**其他补充：**

- 查看安卓系统版本

```bash
adb shell getprop ro.build.version.release
```

- 查看安卓系统 CPU 架构

```bash
adb shell getprop ro.product.cpu.abi
```

### 执行交叉编译

使用刚刚编译好的 `toolchain` 来交叉编译（前提条件：配置环境变量 `NDK_TOOLCHAIN` 为工具链的位置）

```bash
CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC=$NDK_TOOLCHAIN/bin/aarch64-linux-android-gcc go build
```

## 参考文章

- [Go 和 Android 集成实战](https://zhuanlan.zhihu.com/p/86241233)
