[TOC]

# 新系统配置

## mac 2k 外接屏幕配置

为 macOS 开启 HiDPI，让 2K 显示器更舒适

参考链接：https://blog.haitianhome.com/macbook-2k-hidpi.html

## ssh 配置

- 生成 ssh 公钥

命令  `ssh-keygen`，在 ~/.ssh 目录下生成密钥 id_rsa.pub

- ssh 免密输入

https://stackoverflow.com/questions/21095054/ssh-key-still-asking-for-password-and-passphrase

## 配置 golang 环境

### 访问私有的 gitlab 仓库

`go mod tidy` 下载 private repository 出现以下错误。

```bash
	fatal: could not read Username for 'https://git.xxxxx.com': terminal prompts disabled

Confirm the import path was entered correctly.

If this is a private repository, see https://golang.org/doc/faq#git_https for additional information.
```

`go install/mod tidy` 去下载依赖其实是通过 git 命令去下载的，而且默认是 http 协议去下载的，需要读取用户名密码。这里我们修改为 ssh 协议去获取。修改用户目录下的 .gitconfig 文件（没有就创建）

```
➜  ~ cat .gitconfig
# [url "git@git.xxxxx.com:"] 错误
[url "gitlab@git.xxxxx.com:"]
	insteadOf = https://git.xxxxx.com/ # 注意这里的 URL 要和报错的 URL 一致
```

**注意是 gitlab@git.xxxxx.com 不是 git@git.xxxxx.com**

这个前缀看你 ssh clone 的仓库的 URL 来确认。

一开始写成了 `[url "git@git.xxxxx.com:"]`，出现以下错误：

```bash
➜ go mod tidy
	git@git.xxxxx.com: Permission denied (publickey).
	fatal: 无法读取远程仓库。
```

我们可以使用 ssh -T 查看是否能够 ssh 连接访问到该地址

```bash
➜  ~ ssh -T gitlab@git.xxxxx.com
Welcome to GitLab, @xxx!
➜  ~ ssh -T git@git.xxxxx.com 
git@git.xxxxx.com: Permission denied (publickey).
```

如果出现 `dial tcp: lookup git.xxxxx.com on 8.8.8.8:53: no such host`

那么还需要配置一下 `go env -w GOPRIVATE="git.xxxxx.com"`

## golang vscode 配置

因为有些项目又用 vendor 又用 go mod，导致 vscode 无法正确 load packages

1、Open Workspace Settings(JSON)

2、编辑 setting.json 后，关闭 vscode 再重启即可。

```json
{
    "go.toolsEnvVars": {
        "GOFLAGS": "-mod=vendor"
    }
}
```

## 配置 git

> 有些快捷 git 命令安装 oh-my-zh 后就会带上了， ~/.oh-my-zsh/plugins/git/git.plugin.zsh
>
> gp: 推送提交当前到对应的远程分支，alias gp='git push'
>
> gpf: 强制提交，alias gpf='git push --force-with-lease'

具体配置查看之前写的 git 文档

### 解决 git 中文乱码

```bash
git config --global core.quotepath false
```

### 设置 git 的默认编辑器为 vim

Linux 下 git 默认的编辑器使用的是系统默认的文本编辑器。

```bash
git config --global core.editor vim
```

## 安装 brew

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

## 配置 cisco anyconnect vpn

快速打开/关闭 vpn 且不用输入密码。

- 先写好 vpn 的脚本

```bash
➜  xxx cat vpn-connect.sh 
/opt/cisco/anyconnect/bin/vpn -s connect 【VPN Server 地址】 <<"EOF"
【用户名】
【密码】
y
EOF
➜  xxx pwd
/Users/xxx/xxx
```

- 然后在 .zshrc 配置别名

```bash
alias vpnrun="sh /Users/xxx/xxx/vpn-connect.sh"
alias vpnstop="/opt/cisco/anyconnect/bin/vpn -s disconnect 【VPN Server 地址】"
```

- 最后在终端输入 vpnrun 就可以快乐的使用了。

## 安装 python2

因为 python2 新版的 mac 已经启用了，无法通过 `brew install python` 命令简单安装 python2。当前通过 brew 命令安装只能安装 python3。

参考：https://stackoverflow.com/questions/60298514/how-to-reinstall-python2-from-homebrew

### 方法一（推荐）

```bash
Python 2

python, python2 -> python 2.7

# Download/run the legacy macOS installer (pick which one for your sys)
https://www.python.org/downloads/release/python-2716/

# Add pip for python2.7
curl https://bootstrap.pypa.io/pip/2.7/get-pip.py -o get-pip2.py
python2 get-pip2.py

# Optionally check for pip updates (in case of post-eol patches)
python2 -m pip install --upgrade pip

# Optionally add the helpers like easy_install back onto your path
# In your ~/.zprofile or whatever bash/shell profile equivalent
PATH="/Library/Frameworks/Python.framework/Versions/2.7/bin:${PATH}"
export PATH

# Optionally add some helpers while editing shell profile
alias pip2="python2 -m pip"
alias venv2="virtualenv -p python2"
alias venv3="virtualenv -p python3"

# Optionally some apple-specific std libraries are missing, search
# and download them. Example: plistlib.py
curl https://raw.githubusercontent.com/python/cpython/2.7/Lib/plistlib.py -o /Library/Frameworks/Python.framework/Versions/2.7/lib/python2.7/plistlib.py

# Lastly, there is no symlink /usr/bin/python anymore
# /usr/bin is system protected so you can't add one either
# 
# Change your programs to use /usr/local/bin/python
# or google how to disable macOS SIP to make a symlink in /usr/bin
```

这里安装的是 2.7.16 版本。

一开始没用这个步骤的时候，我下载的是 2.7.7 版本。安装的时候显示：安装程序无法安装软件。

后面找了网上的方法说要关闭 SIP，但是关闭了 SIP 也不行。

### 方法二

安装 pyenv，再用 pyenv 安装 python

```bash
brew install pyenv
pyenv install 2.7.18
```

Export PATH if necessary.（重启终端后失效）

```bash
export PATH="$(pyenv root)/shims:${PATH}"
```

Add if necessary.（永久有效）

```bash
echo 'PATH=$(pyenv root)/shims:$PATH' >> ~/.zshrc
```

## 安装其他

- 安装 wget

```bash
brew intall wget
```

- 下载 apipost

https://www.apipost.cn/

- 安装翻译软件 Bob 社区版

https://v0.bobtranslate.com/#/general/quickstart/install

- DiffMerge

可视化的文件比较（也可进行目录比较）与合并工具

http://sourcegear.com/diffmerge/downloads.html

- iPaste

剪切板管理

- 超级右键

为你的右键增加功能，如右键打开终端等等。

- iShot

截图

- [Sequel Ace](https://sequel-ace.com/)

一个 MySQL 数据库管理软件

## vscode 配置

- vscode 保存文件时自动删除行尾空格

点击 setting，搜索 `files.trimTrailingWhitespace`，将选项勾选。

- C++ 设置

[检测到 #include 错误，请更新 includePath](http://t.csdn.cn/eciVI)

## mac 配置

- 将键盘 F1、F2 等键用作标准功能键

打开系统设置，选择键盘，进入键盘的设置窗口，点击**键盘**标签。剩下的懂的都懂。

## 配置 Vim

一键安装配置 vim：https://github.com/chxuan/vimplus

默认的主题看不清楚行号，可以通过下面的方法修改：https://github.com/chxuan/vimplus/issues/91

## 配置 typroa

### Typora + PicGo + Gitee 图床

参考：https://www.modb.pro/db/98359

> 需要注意填写 repo 时候的名称

![gitee-repo](https://gitee.com/goat_cr7/image/raw/master/gitee-repo.png)

ps：gitee 上传的图片太大（超过 1M）导致 typroa 出现 image load failed
