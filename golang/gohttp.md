[TOC]

# Go Http

## 常用函数介绍

- `func (r *Request) FormValue(key string) string`

  POST and PUT body parameters take precedence over URL query string values。翻译一下就是 post 和 put 方法优先去 body 中找对应的字段，如果找不到，就解析 URL 找相应的字段。

- [http 请求头 Range 介绍](https://blog.csdn.net/lv18092081172/article/details/51457525?spm=1001.2101.3001.6650.11&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-11.pc_relevant_paycolumn_v3&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7EBlogCommendFromBaidu%7ERate-11.pc_relevant_paycolumn_v3&utm_relevant_index=17)

- http.ServeContent 会根据请求头中的 Range 来决定文件的响应（传输）范围（一般用来做断点续传或者多线程下载）。

- [multipart/form-data 介绍](https://www.imooc.com/article/314402)
