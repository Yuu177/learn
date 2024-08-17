[TOC]

# FFmpeg

https://github.com/FFmpeg/FFmpeg

## Cross Compiling for Windows

https://trac.ffmpeg.org/wiki/CompilationGuide/CrossCompilingForWindows

```bash
./configure --prefix=./build/windows --arch=x86_64 --target-os=mingw32 --cross-prefix=x86_64-w64-mingw32-
make
make install
```

- nasm/yasm not found or too old. Use --disable-x86asm for a crippled build

`nasm` 和 `yasm` 是用于汇编代码的工具，FFmpeg在编译时需要它们来处理一些性能关键的汇编代码。

```bash
sudo apt-get update
sudo apt-get install nasm yasm
```

### 链接 windows 库

https://github.com/microsoft/vcpkg/pull/14082

```cmake
target_link_libraries(FFmpegDemo
	# ffmepg 库
    swscale
    avformat
    avcodec
    swresample
    avutil
    avdevice
    avfilter

	# 其他库
    pthread
    z # zlib

	# windows 库
    ws2_32 secur32 bcrypt strmiids mfplat mfuuid
)
```

使用 mingw 交叉编译无法链接到对应的 zlib 库，需要安装 mingw-w64 的 zlib

```bash
sudo apt install libz-mingw-w64-dev
```

