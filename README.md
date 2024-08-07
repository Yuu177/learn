# Welcome

之前折腾过买服务器买域名来搭建个人博客，但是发现写了很多的随笔内容，如果以博客的方式发布感觉怪怪的，而且好些内容回过头来看也经常修改。还有就是很多博客的模板是使用标签来分类，这样子文章内容一多就感觉不好整理。

最后找到了最合适自己的记录笔记的方式：**像写代码一样，用树形结构写文章**。所以采用了 GitHub 的 repo 来管理。

## 文档统计

统计日期：2024/01/22

```bash
cloc . --include-lang=Markdown
     150 text files.
     150 unique files.                                          
      13 files ignored.

github.com/AlDanial/cloc v 1.82  T=0.04 s (3112.1 files/s, 513812.8 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Markdown                       138           7343              0          15441
-------------------------------------------------------------------------------
SUM:                           138           7343              0          15441
-------------------------------------------------------------------------------
```

**字数统计：406734**

```bash
find . -name "*.md" -type f -exec cat {} + > ../combined.md
```

合并所有文档后再用 Typora 打开 `combined.md` 查看字数。

## 仓库内容

个人总结的学习笔记、随笔

## 笔记软件

[Typora 一款 Markdown 编辑器和阅读器](https://typoraio.cn/)

## 排版规范

此仓库遵循 [中文排版指南](https://github.com/sparanoid/chinese-copywriting-guidelines) 规范（TODO）

## ~~约定式提交~~

~~此仓库 commit 遵循 https://www.conventionalcommits.org/zh-hans/v1.0.0/~~
