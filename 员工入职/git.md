[TOC]

# Git

## 新建远程仓库和本地仓库关联

- github

```bash
…or create a new repository on the command line
 
echo "# test" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:Yuu177/test.git
git push -u origin main
 
…or push an existing repository from the command line
 
git remote add origin git@github.com:Yuu177/test.git
git branch -M main
git push -u origin main
```

- gitlab

```bash
git remote rename origin old-origin
git remote add origin gitlab@git.xxxx.com:panyu.tan/test.git
git push -u origin --all
git push -u origin --tags
```

- 对以上命令补充

```bash
git remote add [shortname] [url] # shortname 为远程库的名字，对应一个 url
git remote rm name  # 删除远程仓库
git remote rename old_name new_name  # 修改本地仓库名

# -M: 是 --move --force 的缩写。
# --move(-m): Move/rename a branch.
# --force(-f): 即使新命名的 branch 名存在也执行。
git branch -M main

# 由于远程库是空的，我们第一次推送 master 分支时，加上了 -u 参数，Git 不但会把本地的 master 分支内容推送的远程新的 master 分支，还会把本地的 master 分支和远程的 master 分支关联起来，在以后的推送或者拉取时就可以简化命令。
git push -u origin main
```

## 回滚

- 本地代码回滚到上一版本（或者指定版本）

`git reset --hard HEAD^`

- 完成撤销,同时将代码恢复到前一 commit_id 对应的版本。

`git reset --hard commit_id`

- 完成 commit 命令的撤销，但是不对代码修改进行撤销，可以直接通过 git commit 重新提交对本地代码的修改。

`git reset commit_id`

## 新建分支并推送

1. 新建分支

`git checkout -b branch_name`

2. 推送本地新建的分支到远程仓库

`git push --set-upstream origin branch_name`

ps：推送本地新建的分支到远程仓库（此时远程仓库还没有对应的分支，--set-upstream 的作用就在此）。输入命令后会显示 'branch_name' 设置为跟踪来自 'origin' 的远程分支 'branch_name'。

3. 强制推送本地分支覆盖掉远程分支

`git push -f`

ps：强制推送本地分支覆盖掉远程分支。只适合自己 create 的远程分支上进行。这样子远程分支上就不会有多个 commit

- 从指定的 commit 节点创建分支

`git checkout -b branchName commit_id`

## 重命名分支&删除分支

1. 本地分支重命名

`git branch -m oldName newName`

2. 将重命名后的分支推送到远程

`git push origin newName`

3. 删除远程的旧分支

`git push --delete origin oldName`

4. 删除本地分支

`git branch -D branchName`

5、您的分支基于 'origin/oldName'，但此上游分支已经不存在。
  （使用 "git branch --unset-upstream" 来修复）

6、为推送当前分支并建立与远程上游的跟踪，使用

git push --set-upstream origin newName

## 分支同步

- 将 master 分支的代码（远程最新提交）同步到 checkout 出来的分支

`git rebase master`

## 本地分支对应多个不同仓库的远程分支

假设本地分支 main。对应有 gitlab 的仓库 gitlab/main 和 github 的仓库 github/main。我们如何设置默认 git pull 或者 git push 默认操作的远程仓库？

--set-upstream 为 git pull/fetch 设置上游（--set-upstream 远程仓库名 远程分支名）

`git pull --set-upstream origin main`

## 合并多个 commit 为一个 commit

- 当前 git log 记录

```bash
* fa7fd6f22 - five ci (2 分钟前) <Panyu Tan>
* 28692885b - four ci (3 分钟前) <Panyu Tan>
* a305a466b - three ci (31 分钟前) <Panyu Tan>
* 72e242e40 - second ci (31 分钟前) <Panyu Tan>
* bcbb966e2 - first ci (32 分钟前) <Panyu Tan>
* 6a97ce7d2 - fix public fiels content (20 小时前) <Panyu Tan>
```

### 场景一

我们要把 first ~ five 五个提交合并为一个 merge all

- `git rebase -i 6a97ce7d2`

从 `6a97ce7d2 - fix public fiels content` **下一个提交**到最新的节点进行合并。这个命令就是把 first ci 的提交到最新的提交(five ci) 准备要进行合并 commit。

其中 -i 的意思是 --interactive，即弹出交互式的界面让用户编辑完成合并操作。

接下来会进入 vi 界面。如下：

```bash
pick bcbb966e2 first ci
pick 72e242e40 second ci                                                                                                                                                       
pick a305a466b three ci
pick 28692885b four ci
pick fa7fd6f22 five ci

# p, pick <提交> = 使用提交
# s, squash <提交> = 使用提交，但挤压到前一个提交
```

编辑修改为：

```bash
pick bcbb966e2 first ci  # 使用 first 提交记录
s    72e242e40 second ci # second 合并到 first                                                                                                                                               
s    a305a466b three ci  # three 合并到上一个提交，因为 second 合并到 first 了，所以 three 也会合并到 first
s    28692885b four ci   # 同上
s    fa7fd6f22 five ci   # 同上
```

修改完后 wq 保存。接下来再次进入 vi 的编辑界面，重新编辑提交记录。

```bash
# 这是一个 5 个提交的组合。                                                                                                                                                      
# 这是第一个提交说明：

first ci
 
# 这是提交说明 #2：
 
second ci
 
# 这是提交说明 #3：
  
three ci

# 这是提交说明 #4：
 
four ci

# 这是提交说明 #5：
 
five ci
```

因为 # 行会被忽略，所以我们只保留一行没有 # 即可。如下

```bash
# 这是一个 5 个提交的组合。                                                
merge all
```

wq 保存退出。接下我们再来查看 git log。发现 5 个提交已经被合并为一个了。

```bash
* 60cc96e41 - merge all (3 分钟前) <Panyu Tan>
* 6a97ce7d2 - fix public fiels content (20 小时前) <Panyu Tan>
```

### 场景二

我想把 first 和 second 两个合并为一个commit。three four five 这三个提交合并为一个 commit。

根据刚才的实践，我们只需要这样子修改即可。

```bash
pick bcbb966e2 first ci  # 使用 first 提交记录
s    72e242e40 second ci # second 合并到 first                                                                                                                                               
pick a305a466b three ci  # 使用 three 提交记录
s    28692885b four ci   # four 合并到 three 
s    fa7fd6f22 five ci   # five 合并到 three
```

## git 问题

- git status 乱码

参考链接：https://blog.csdn.net/u012145252/article/details/81775362

- git pull 的时候出现错误

> kex_exchange_identification: read: Connection reset by peer
> fatal: Could not read from remote repository.
>
> Please make sure you have the correct access rights
> and the repository exists.

设置一下 git config http.sslVerify "false"

参考链接：https://stackoverflow.com/questions/54611871/ssh-exchange-identification-read-connection-reset-by-peer-error-when-trying

## git 快捷别名配置

安装完 Oh-my-zsh 会自动配置 git 命令的别名: ~/.oh-my-zsh/plugins/git/git.plugin.zsh
