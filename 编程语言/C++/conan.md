[TOC]

# conan

Conan 是一个可以帮 C/C++ 进行依赖管理的包管理器。

参考官方教程：https://docs.conan.io/en/latest/getting_started.html

## 问题

### 自动下载安装依赖

```bash
xorg/system: ERROR: while executing system_requirements(): System requirements: 'libx11-xcb-dev, libfontenc-dev, libxaw7-dev, libxkbfile-dev, libxmu-dev, libxmuu-dev, libxpm-dev, libxres-dev, libxss-dev, libxtst-dev, libxv-dev, libxvmc-dev, libxxf86vm-dev, libxcb-render-util0-dev, libxcb-xkb-dev, libxcb-icccm4-dev, libxcb-image0-dev, libxcb-keysyms1-dev, libxcb-randr0-dev, libxcb-shape0-dev, libxcb-sync-dev, libxcb-xfixes0-dev, libxcb-xinerama0-dev, libxcb-dri3-dev' are missing but can't install because tools.system.package_manager:mode is 'check'.Please update packages manually or set 'tools.system.package_manager:mode' to 'install' in the [conf] section of the profile, or in the command line using '-c tools.system.package_manager:mode=install'
ERROR: Error in system requirements
```

编译 OpenCV 相关代码的时候出现以上问题：缺少相关依赖。

我们可以通过 apt install 来一个一个安装依赖。

```
sudo apt install libx11-xcb-dev
```

或者配置 conan 来自动安装。配置 `~/.conan/global.conf`，没有就新建。

```bash
➜  .conan cat global.conf 
tools.system.package_manager:mode = install
tools.system.package_manager:sudo = True
```

### 设置 C++11 标准库

- 静态链接失败

```bash
[100%] Linking CXX executable bin/CppDemo
/usr/bin/ld: CMakeFiles/CppDemo.dir/main.cpp.o: in function `main':
main.cpp:(.text.startup+0x98): undefined reference to `cv::imshow(std::__cxx11::basic_string<char, std::char_traits<char>, std::allocator<char> > const&, cv::_InputArray const&)'
collect2: error: ld returned 1 exit status
make[2]: *** [CMakeFiles/CppDemo.dir/build.make:97：bin/CppDemo] 错误 1
make[1]: *** [CMakeFiles/Makefile2:83：CMakeFiles/CppDemo.dir/all] 错误 2
make: *** [Makefile:91：all] 错误 2
```

- Conan profile（conan install 的时候会打印输出这些配置信息）

```bash
Configuration:
[settings]
arch=x86_64
arch_build=x86_64
build_type=Release
compiler=gcc
compiler.libcxx=libstdc++
compiler.version=9
os=Linux
os_build=Linux
```

- 解决

I guess you don't honor the CXX ABI of your conan profile when you build your project:

```bash
undefined reference to `cv::imread(std::__cxx11::basic_string<char
```

I see `std::__cxx11::basic_string<char`, and you have `libstdc++` in your profile, instead of `libstdc++11`

所以我们可以修改 conan install 命令，添加 `-s compiler.libcxx=libstdc++11`

```bash
conan install ../conanfile.txt --build=missing
# 修改为
conan install ../conanfile.txt --build=missing -s compiler.libcxx=libstdc++11
```

- 参考

https://github.com/conan-io/conan-center-index/issues/10448