[TOC]

## docker

### 安装

`brew install --cask --appdir=/Applications docker`

### 常用命令

```bash
# 查看容器
docker ps -a
# 停止
docker stop id
# 删除
docker rm id

docker run：根据镜像创建一个容器并运行一个命令，操作的对象是镜像；
docker exec：在运行的容器中执行命令，操作的对象是 容器。
-p: 指定端口映射，格式为：主机(宿主)端口:容器端口
-i：以交互模式运行容器
-t：为容器重新分配一个伪输入终端
-d：后台运行容器，并返回容器 ID
-e username="ritchie": 设置环境变量；
--name="nginx-lb": 为容器指定一个名称。
```

## 参考链接

https://www.runoob.com/docker/docker-tutorial.html