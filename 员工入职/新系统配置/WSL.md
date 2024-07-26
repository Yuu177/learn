[TOC]

# WSL

## 检测到 localhost 代理配置

开启代理打开 WSL 后出现在终端第一行出现如下的信息提示：

```bash
wsl: 检测到 localhost 代理配置，但未镜像到 WSL 。NAT 模式下的 WSL 不支持 localhost 代理
```

- 解决方法

在资源管理器输入 `%userprofile%`，在打开的文件夹下新建 `.wslconfig` 文件，内容如下：

```bash
[experimental]
autoMemoryReclaim=gradual # 开启自动回收内存，可在 gradual, dropcache, disabled 之间选择
networkingMode=mirrored # 开启镜像网络
dnsTunneling=true # 开启 DNS Tunneling
firewall=true # 开启 Windows 防火墙
autoProxy=true # 开启自动同步代理
sparseVhd=true # 开启自动释放 WSL2 虚拟硬盘空间
```

在 windows 上使用 `wsl --shutdown` 命令结束 wsl ，再重新运行 wsl 就可以了。

- 参考

https://www.v2ex.com/t/992252