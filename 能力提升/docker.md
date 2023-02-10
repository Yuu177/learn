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
docker ps -a：显示全部容器
docker ps：显示当前运行的容器
docker stop 容器ID：停止容器
docker rm 容器ID：删除容器
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

## docker 和本机文件互相拷贝

### docker cp

- 从 docker 容器拷贝文件到本机

```
docker cp container_created:path <path>
```

- 将文件从本机拷贝到 docker 容器

```
docker cp <path> container_created:path
```

- 例子

```
docker cp hello.py f78c63afeb86:/home
```

命令说明：拷贝本地文件 `hello.py` 到容器 `f78c63afeb86` 的 `/home` 目录下。

### 挂载/共享文件夹

当然也可以通过命令参数 `-v` 来共享文件夹。

```
-v 本地文件夹路径:docker路径
```

## docker 挂载摄像头并显示图片

- 问题

1、在 docker 容器运行 OpenCV 显示图片的时候出现错误：cannot open display

2、调用摄像头的时候显示：can't open camera by index 0

> 以下命令都是在本地主机执行

- 安装相关工具并执行命令

```bash
sudo apt-get install x11-xserver-utils
xhost +
```

Linux 系统目前的主流图像界面服务 `X11` 支持客户端/服务端（C/S）的工作模式，只要在容器启动的时候，将 『unix:端口』或『主机名:端口』共享给 Docker，Docker 就可以通过端口找到显示输出的地方，和 Linux 系统共用显示接口。

`xhost +` 命令的作用是开放权限，允许所有用户访问显示接口。也可以指定特定用户：

```bash
xhost +local:docker # 只允许 Docker 用户访问显示接口 
```

注意：每次重新开机，都需要再次执行 `xhost +`。

- 执行 docker 命令

```bash
docker run -t -i \
-v /etc/localtime:/etc/localtime:ro \
-v /tmp/.X11-unix:/tmp/.X11-unix \
-e DISPLAY=unix$DISPLAY \
--device=/dev/video0 \
--device=/dev/video1 \
--device=/dev/video2 \
--device=/dev/video3 \
-v /home/jinx/code:/root/code \
ubuntu18.04_dev_v1_test
```

- 命令解释

```bash
docker run -t -i \ # -t 为容器重新分配一个伪输入终端，-i 以交互模式运行容器
-v /etc/localtime:/etc/localtime:ro \ #（可选）-v 共享/挂载目录。docker 容器时间同步，ro 代表只读属性
-v /tmp/.X11-unix:/tmp/.X11-unix \ # 共享本地 unix 端口
-e DISPLAY=unix$DISPLAY \ # -e 设置环境变量。修改环境变量 DISPLAY
--device=/dev/video0 \ # 添加主机设备给容器，相当于设备直通
--device=/dev/video1 \
--device=/dev/video2 \
--device=/dev/video3 \
-v /home/jinx/code:/root/code \ # 共享本地代码目录
ubuntu18.04_dev_v1_test # 镜像名称
```

使用 `ls /dev/ | grep video*` 查看系统摄像头，然后把它们全都挂载上去。

```bash
➜  ~ ls /dev/ | grep video*
video0
video1
video2
video3
```

## 镜像拉取、保存、打包和加载

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
- [Docker 挂载摄像头并显示图像](https://blog.csdn.net/weixin_40922744/article/details/103245166)