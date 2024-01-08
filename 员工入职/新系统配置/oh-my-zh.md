[TOC]

# oh-my-zh

## 安装

https://github.com/ohmyzsh/ohmyzsh

修改默认 shell 为 zsh（需要重启终端 or 电脑）

```bash
chsh -s /bin/zsh
```

## 插件

### zsh-autosuggestions & zsh-syntax-highlighting

历史命令（zsh-autosuggestions），命令高亮（zsh-syntax-highlighting）

```bash
git clone https://github.com/zsh-users/zsh-autosuggestions
git clone https://github.com/zsh-users/zsh-syntax-highlighting
```

下载好后移动到 `$ZSH/plugins/` 目录下：

```bash
sudo mv zsh-autosuggestions zsh-syntax-highlighting $ZSH/plugins/
```

在文件 `~/.zshrc` 最后写入：

```bash
source $ZSH/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh
source $ZSH/plugins/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
```

保存后终端输入 `source ~/.zshrc` 并重启终端。

### autojump

https://github.com/wting/autojump

- docker 安装错误

docker 中安装 autojump 出现 Unsupported shell: None，`echo $SHELL` 发现是空

1、使用 zsh 打开 docker：`docker exec -it 40dfecca1f53 /bin/zsh`（这里不能用 `/bin/bash` 打开）

2、执行命令：`export SHELL=/bin/zsh`，重新执行安装步骤即可

## docker zsh

- 报错：zsh (anon):12: character not in range

在 `.zshrc` 下添加

```bash
export LC_ALL=C.UTF-8
export LANG=C.UTF-8
```

- 显示用户和计算机名

修改 `.zshrc` 的 `ZSH_THEME="agnoster"`

- 退出 docker 后再登陆时，zsh 配置不生效


每次登陆后需要 `source ~/.zshrc`，或者在 `/etc/zsh/zshrc` 文件末尾增加 `source ~/.zshrc`

