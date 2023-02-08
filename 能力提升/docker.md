[TOC]

# docker

## 安装

- mac

```bash
brew install --cask --appdir=/Applications docker
```

- Linux

```bash
sudo apt update && sudo apt install docker.io
```

安装完后执行 `docker ps` 显示用户没有权限。参考：[Docker 启动 Get Permission Denied ](https://www.cnblogs.com/informatics/p/8276172.html)，可能需要重启电脑。

## 常用命令

参考：[Docker 常用命令](https://www.w3cschool.cn/docker/docker-nx3g2gxn.html)

- 命令

```bash
docker ps -a：查看容器
docker stop id：停止容器
docker rm id：删除容器
docker run：根据镜像创建一个容器并运行一个命令，操作的对象是镜像
docker exec：在运行的容器中执行命令，操作的对象是容器
```

- 参数

```bash
-p: 指定端口映射，格式为：主机端口:容器端口
-i：以交互模式运行容器
-t：为容器重新分配一个伪输入终端
-d：后台运行容器，并返回容器 ID
-v：给容器挂载存储卷，挂载到容器的某个目录。格式为：主机目录:容器目录
-e username="hello": 设置环境变量
--name="world": 为容器指定一个名称
```

## docker 显示 GUI

在 docker 容器运行 OpenCV 显示图片的时候出现错误：cannot open display

> 以下命令都是在本地主机执行

- 安装相关工具

```bash
sudo apt-get install x11-xserver-utils
xhost +
```

这两句命令的作用是开放权限，允许所有用户（当然包括 docker）访问 X11 的显示接口。每次重新开机，需要在本机操作一次 `xhost +`。

- 执行 docker 命令

```bash
docker run -t -i \
-v /etc/localtime:/etc/localtime:ro \
-v /tmp/.X11-unix:/tmp/.X11-unix \
-e DISPLAY=unix$DISPLAY \
-e GDK_SCALE \
-e GDK_DPI_SCALE \
ubuntu18.04_dev_v1
```

- 命令解释

```bash
docker run -t -i \ # -t 为容器重新分配一个伪输入终端，-i 以交互模式运行容器
-v /etc/localtime:/etc/localtime:ro \ # -v 共享目录。docker 容器时间同步，ro 代表只读属性
-v /tmp/.X11-unix:/tmp/.X11-unix \ # 共享本地 unix 端口
-e DISPLAY=unix$DISPLAY \ # -e 设置环境变量。修改环境变量 DISPLAY
-e GDK_SCALE \ # 未知
-e GDK_DPI_SCALE \ # 未知
ubuntu18.04_dev_v1 # 镜像
```

## 镜像保存、打包和加载

- 拉取 ubuntu18.04 镜像

```
docker pull ubuntu:18.04
```

- 进入 ubuntu:18.04 容器

```
docker run -t -i ubuntu:18.04
```

安装开发环境

```bash
root@49708f314238:/# apt update
root@49708f314238:/# apt install g++ cmake git
root@49708f314238:/# exit
```

将上面修改过的容器保存成新镜像 ubuntu_dev

```
docker commit -m 'install dev env' -a="zhangsan" 49708f314238 ubuntu_dev
```

`49708f314238` 为上面运行镜像的容器 ID

运行新镜像，就可以在里面 clone 代码进行编译了。

```
docker run -i -t ubuntu_dev
```

- 镜像打包

有时候我们需要把镜像共享给别人。

使用 save 命令指定打包后的存放的路径

```bash
docker save ubuntu_dev > ~/ubuntu_dev.tar
```

打包完成可以在相应目录下看到多了一个 tar 包，这就是装好环境的 docker 镜像。

- 镜像导入

```bash
docker load -i ubuntu_dev.tar
```

## 参考链接

- [Docker 入门实战](https://www.w3cschool.cn/docker/docker-tutorial.html)