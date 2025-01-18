[TOC]

# typecho 网站搭建

## LNMP 环境安装

[手动搭建 LNMP 环境（CentOS 8）](https://cloud.tencent.com/document/product/213/49304)

因为安装 typecho 需要依赖 LNMP 环境，所以无脑按照腾讯云的步骤来就可以了。以下步骤是建立在按照这个基础上修改的，请务必对齐配置。

## typecho 部署

### 上传  typecho.zip 文件

在 linux 机器上创建好网站目录，并上传下载好的 [typecho.zip](http://typecho.org/download) 文件。不建议创建目录在 /root 下面。我试过发现怎么也访问不到，后面换其他目录就好了。

这里我们选择上传文件到 /usr/share/nginx/html/ 下。解压 typecho.zip 后再添加权限 `chmod -R 777 html`。因为待会安装需要网站目录的读写权限。

### 配置 nginx

> 配置之前记得备份一下文件，下面的配置文件都是默认开启 https。

- /etc/nginx/conf.d/default.conf
- /etc/nginx/nginx.conf

#### 修改 default.conf

配置文件照抄就行。

>ssl 的证书文件和私钥文件上传到 /etc/nginx/ 目录下用来开启 https，参考后续内容。

```bash
# 如果不开启 https，那么把 ssl_ 前缀的相关行注释掉即可。并且把 listen 433 ssl 改为 listen 80
server {
    listen 443 ssl; # 监听 433 ssl 端口
    server_name  cr7.life; # 填写你的网站域名
    
    # ssl 配置
    ssl_certificate cr7.life_bundle.crt; # 填写您的证书文件名称
    ssl_certificate_key cr7.life.key; # 填写您的私钥文件名称
    ssl_session_timeout 5m;
    ssl_protocols  TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;
    ssl_prefer_server_ciphers on;

    root /usr/share/nginx/html; # 你的网站的目录

    error_page   500 502 503 504  /50x.html;
    location = /50x.html {
        root   /usr/share/nginx/html;
    }

    include /etc/nginx/default.d/*.conf;

    location ~ \.php$ {
        include        fastcgi_params;
        fastcgi_pass unix:/run/php-fpm/www.sock;
	    client_max_body_size 20m;
        fastcgi_connect_timeout 30s;
        fastcgi_send_timeout 30s;
        fastcgi_read_timeout 30s;
        fastcgi_intercept_errors on;
    }
}

# 如果不开启 https，就把这个 server 去掉
server {
    listen 80;
    server_name cr7.life;    # 填写您的证书绑定的域名
    return 301 https://$host$request_uri; # 将 http 的域名请求转成 https
}
```

#### 修改 nginx.conf

配置文件照抄。

```bash
user  nginx;
worker_processes auto;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

include /usr/share/nginx/modules/*.conf;

events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile            on;
    tcp_nopush          on;
    tcp_nodelay         on;
    keepalive_timeout   65;
    types_hash_max_size 2048;

    #gzip  on;
	
	# https ssl 配置。如果不开启 https，注释掉即可。
    ssl_certificate cr7.life_bundle.crt; # 填写您的证书文件名称
    ssl_certificate_key cr7.life.key; # 填写您的私钥文件名称

    include /etc/nginx/conf.d/*.conf;
}
```

修改完配置文件后，输入命令 `nginx -t ` 检查配置文件是否有问题。

重启一下 nginx 使配置生效 `sudo systemctl restart nginx`。

### 创建 mysql 数据库

登录到 mysql 创建待会需要用到的据库名，如：

```sql
create database my_blog;
```

### 开始安装 typecho

在浏览器输入 `ip/install.php` 进入安装步骤。一开始我们创建的数据库名这里就派上用场了。填入完信息，大功告成。

***注意**：**站点信息**这一栏一开始我们没有申请域名的话可以先填你的外网 ip。后面如果要用域名访问的话，需要修改成你的域名地址。*

### 域名与 DNS 解析设置

可以给自己的 Typecho 网站设定一个单独的域名。别人可以使用易记的域名访问您的网站，而不需要使用复杂的 IP 地址。下面是腾讯云的操作步骤：

1. 通过腾讯云 [购买域名](https://dnspod.cloud.tencent.com/?from=qcloud)，具体操作请参考 [域名注册](https://cloud.tencent.com/document/product/242/9595)。
2. 进行 [网站备案](https://cloud.tencent.com/product/ba?from=qcloudHpHeaderBa&fromSource=qcloudHpHeaderBa)。
   域名指向中国境内服务器的网站，必须进行网站备案。在域名获得备案号之前，网站是无法开通使用的。您可以通过腾讯云免费进行备案，审核时长请参考 [备案审核](https://cloud.tencent.com/document/product/243/19650)。
3. 通过腾讯云 [DNS解析 DNSPod](https://cloud.tencent.com/product/cns?from=qcloudHpHeaderCns&fromSource=qcloudHpHeaderCns) 配置域名解析。具体操作请参考 [A 记录](https://cloud.tencent.com/document/product/302/3449)，将域名指向一个 IP 地址（外网地址）。

### 开启 HTTPS 访问

可参考 [安装 SSL 证书](https://cloud.tencent.com/document/product/1207/47027) 文档为您的 Typecho 实例安装 SSL 证书并开启 HTTPS 访问。

## 补充其他

### PHP-FPM

PHP-FPM 负责管理一个进程池来处理来自 Web 服务器的 HTTP 动态请求，在 PHP-FPM 中，master 进程负责与 Web 服务器进行通信，接收 HTTP 请求，再将请求转发给 worker 进程进行处理，worker 进程主要负责动态执行 PHP 代码，处理完成后，将处理结果返回给 Web 服务器，再由 Web 服务器将结果发送给客户端。这就是 PHP-FPM 的基本工作原理。

