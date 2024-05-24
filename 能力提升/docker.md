[TOC]

# docker

可以把 docker 简单理解为一个小型虚拟机。

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

## 镜像拉取、保存、打包和加载

- 拉取 ubuntu18.04 镜像

```bash
docker pull ubuntu:18.04
```

- 通过 ubuntu:18.04 镜像创建一个容器并且运行

```bash
docker run -it ubuntu:18.04 /bin/bash
```

- 安装开发环境

```bash
root@49708f314238:/# apt update
root@49708f314238:/# apt install g++ cmake git
root@49708f314238:/# exit
```

- 将上面修改过的容器保存成新镜像 ubuntu_dev

```bash
docker commit -m 'install dev env' -a="zhangsan" 49708f314238 ubuntu_dev
```

> `49708f314238` 为上面运行镜像的容器 ID

- 镜像打包

有时候我们需要把镜像共享给别人。使用 save 命令指定打包后的存放的路径

```bash
docker save ubuntu_dev > ~/ubuntu_dev.tar
```

打包完成可以在相应目录下看到多了一个 tar 包，这就是装好环境的 docker 镜像。

- 镜像导入

```bash
docker load -i ubuntu_dev.tar
```

## docker 命令说明

```bash
docker run -it ubuntu /bin/bash
```

- `run`：根据镜像创建一个容器并运行一个命令
- `-i`：交互式操作
- `-t`：终端
- `ubuntu`：镜像名称
- `/bin/bash`：放在镜像名后的是命令，这里我们希望有个交互式 Shell，因此用的是 `/bin/bash`

通过 `docker ps` 我们可以看到刚才我们启动的容器信息。注意容器和镜像的区别，用代码来做比较的话，**镜像就相当于类，而容器相当于该类的实例**。

> docker ps 的输出

| CONTAINER ID | IMAGE  | COMMAND     | CREATED       | STATUS       | PORTS | NAMES          |
| ------------ | ------ | ----------- | ------------- | ------------ | ----- | -------------- |
| 40dfecca1f53 | ubuntu | "/bin/bash" | 5 minutes ago | Up 5 minutes |       | peaceful_noyce |

上面 `run` 命令是通过一个镜像创建一个容器，容器创建好后，只要我们不删除该容器，里面的数据就会一直保存。

如果需要再次进入该容器的话（前提是该容器处于运行态，如果不是则通过 `docker start <CONTAINER ID>` 运行容器），执行 `exec` 命令即可：

```bash
docker exec -it 40dfecca1f53 /bin/bash
```

`run` 的操作对象是镜像，`exec` 操作对象是容器。**不需要每次都 `run` 创建一个新的容器**，这是初学者容易犯的一个错误。

### 常用命令

```bash
docker ps -a：显示全部状态的容器
docker ps：显示当前运行（运行态）的容器
docker stop <容器ID>：停止容器
docker rm <容器ID>：删除容器（注意，容器里的数据会被删除）
docker start <容器ID>：启动容器
docker run：根据镜像创建一个容器并运行一个命令，操作的对象是镜像
docker exec <容器ID>：在运行状态的容器中执行命令，操作的对象是容器
docker rmi：删除镜像
```

### 常用参数

```bash
-p：指定端口映射，格式为：<主机端口>:<容器端口>
-i：以交互模式运行容器
-t：为容器重新分配一个伪输入终端
-d：后台运行容器，并返回容器 ID
-v：给容器挂载存储卷，挂载到容器的某个目录。格式为：<主机目录>:<容器目录>
-e username="hello"：设置环境变量
--name="helloworld"：为容器指定一个名称
-u：指定用户，uid:gid
```

其他参考：[Docker 常用命令](https://www.w3cschool.cn/docker/docker-nx3g2gxn.html)

## 容器的状态

容器的五种状态：

```bash
created：初建状态
running：运行状态
exited（stopped）：停止状态
paused： 暂停状态
deleted：删除状态
```

容器在执行某种命令后进入的状态：

```bash
docker create：创建容器后，不立即启动运行，容器进入初建状态；
docker run：创建容器，并立即启动运行，进入运行状态；
docker start：容器转为运行状态；
docker stop：容器将转入停止状态；
docker kill：容器在故障（死机）时，执行 kill（断电），容器转入停止状态，这种操作容易丢失数据，除非必要，否则不建议使用；
docker restart：重启容器，容器转入运行状态；
docker pause：容器进入暂停状态；
docker unpause：取消暂停状态，容器进入运行状态；
docker rm：删除容器，容器转入删除状态（如果没有保存相应的数据库，则状态不可见）。
```

## 挂载摄像头并显示图像

1、在 docker 容器运行 OpenCV 显示图片的时候出现错误：cannot open display

2、调用摄像头的时候显示：can't open camera by index 0

### 挂载摄像头

使用 `ls /dev/ | grep video*` 查看系统摄像头，然后把它们全都挂载上去。

```bash
➜  ~ ls /dev/ | grep video*
video0
video1
video2
video3
```

通过 `--device` 参数挂载摄像头，这样子在 docker 中就可以使用宿主机的摄像头了。

```bash
docker run -t -i \
--device=/dev/video0 \
--device=/dev/video1 \
--device=/dev/video2 \
--device=/dev/video3 \
ubuntu18.04
```

### 显示图像

**以下命令都是在本地主机执行**

- 安装相关工具

```bash
sudo apt-get install x11-xserver-utils
```

- 执行命令

```bash
xhost +
```

`xhost +` 命令的作用是开放权限，允许所有用户访问显示接口。也可以指定特定用户：

```bash
xhost +local:docker # 只允许 Docker 用户访问显示接口
```

**注意：每次重新开机，都需要再次执行 `xhost +`。**

Linux 系统目前的主流图像界面服务 `X11` 支持客户端/服务端（C/S）的工作模式，只要在容器启动的时候，将 『unix:端口』或『主机名:端口』共享给 Docker，Docker 就可以通过端口找到显示输出的地方，和 Linux 系统共用显示接口。

```bash
docker run -t -i \
-v /tmp/.X11-unix:/tmp/.X11-unix \ # 共享本地 unix 端口
-e DISPLAY=unix$DISPLAY \ # 修改环境变量 DISPLAY
ubuntu18.04
```

## 避免 docker 挂载时产生 root 权限文件

在 docker 挂载磁盘的时候，由于容器内默认是 root 用户，会导致挂载中产生的文件属于 root:root。而一般容器外用户并不是 root，如果操作在 docker 中生成的共享文件的时候会导致一些权限问题。

解决办法：在 docker 中创建一个 uid 和 gid 和宿主机一致的用户。

1. 使用 id 命令查看宿主机的 uid 和 gid，一般都是 1000。

2. 登录到 docker linux（默认用户为 root）

```bash
docker exec -it 5af149ba1e95 /bin/bash
```

3. 创建和宿主机一样的 id 和 group 的用户，并添加密码

```bash
groupadd -g 1000 mygroup
useradd -u 1000 -g 1000 -m myuser
passwd myuser
```

> `-m`：创建 /home 目录
>
> 删除用户账户：`sudo userdel -r username`。`-r`：删除用户账户以及 home 目录

4. 设置 sudo 权限

```bash
chmod 744 /etc/sudoers
vi /etc/sudoers
# 在文件最后一行添加 myuser ALL=(ALL) NOPASSWD:ALL
chmod 400 /etc/sudoers
```

5. exit 退出该容器
6. 使用 uid 为 1000 和 gid 为 1000 的用户登录

```bash
docker exec -it -u 1000:1000 5af149ba1e95 /bin/bash
```

## 示例

```bash
docker run -t -i \
-v /etc/localtime:/etc/localtime:ro \
-v /tmp/.X11-unix:/tmp/.X11-unix \
-e DISPLAY=unix$DISPLAY \
--device=/dev/video0 \
--device=/dev/video1 \
--device=/dev/video2 \
--device=/dev/video3 \
-v /mnt:/mnt \
-v /home/jinx/code:/home/jinx/code \
ubuntu
```

- 命令解释

```bash
docker run -t -i \ # run 创建容器，-t 为容器重新分配一个伪输入终端，-i 以交互模式运行容器
-v /etc/localtime:/etc/localtime:ro \ #（可选）-v 共享/挂载目录。docker 容器时间同步，ro 代表只读属性
-v /tmp/.X11-unix:/tmp/.X11-unix \ # 共享本地 unix 端口
-e DISPLAY=unix$DISPLAY \ # -e 设置环境变量。修改环境变量 DISPLAY
--device=/dev/video0 \ # 添加主机设备给容器，相当于设备直通
--device=/dev/video1 \
--device=/dev/video2 \
--device=/dev/video3 \
-v /mnt:/mnt \ # 挂载目录
-v /home/jinx/code:/home/jinx/code \ # 保持 docker 中的路径和本地的路径一致
ubuntu # 镜像名称
```

## 网络无法连接问题

```bash
root@f22965c0097d:/# apt update
Err:1 http://mirrors.aliyun.com/ubuntu bionic InRelease
  Temporary failure resolving 'mirrors.aliyun.com'
Err:2 http://mirrors.aliyun.com/ubuntu bionic-security InRelease
  Temporary failure resolving 'mirrors.aliyun.com'
Err:3 http://mirrors.aliyun.com/ubuntu bionic-updates InRelease
  Temporary failure resolving 'mirrors.aliyun.com'
Err:4 http://mirrors.aliyun.com/ubuntu bionic-proposed InRelease
  Temporary failure resolving 'mirrors.aliyun.com'
Err:5 http://mirrors.aliyun.com/ubuntu bionic-backports InRelease
  Temporary failure resolving 'mirrors.aliyun.com'
Reading package lists... Done
Building dependency tree       
Reading state information... Done
8 packages can be upgraded. Run 'apt list --upgradable' to see them.
W: Failed to fetch http://mirrors.aliyun.com/ubuntu/dists/bionic/InRelease  Temporary failure resolving 'mirrors.aliyun.com'
W: Failed to fetch http://mirrors.aliyun.com/ubuntu/dists/bionic-security/InRelease  Temporary failure resolving 'mirrors.aliyun.com'
W: Failed to fetch http://mirrors.aliyun.com/ubuntu/dists/bionic-updates/InRelease  Temporary failure resolving 'mirrors.aliyun.com'
W: Failed to fetch http://mirrors.aliyun.com/ubuntu/dists/bionic-proposed/InRelease  Temporary failure resolving 'mirrors.aliyun.com'
W: Failed to fetch http://mirrors.aliyun.com/ubuntu/dists/bionic-backports/InRelease  Temporary failure resolving 'mirrors.aliyun.com'
W: Some index files failed to download. They have been ignored, or old ones used instead.
```

使用 host 模式解决：https://stackoverflow.com/a/66714888/24490421

网络模式说明：[Docker 网络模式](https://blog.51cto.com/u_16099316/6467190)

## 参考链接

- [Docker 入门实战](https://www.w3cschool.cn/docker/docker-tutorial.html)
- [Docker 挂载摄像头并显示图像](https://blog.csdn.net/weixin_40922744/article/details/103245166)