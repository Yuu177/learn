[TOC]

# 深度学习

## 数据格式

深度学习框架中，数据一般是 4D，用 NCHW 或 NHWC 表达，其中：

```bash
N - Batch
C - Channel
H - Height
W - Width
```

### RGB 图像数据举例

表达 RGB 彩色图像时，一个像素的 RGB 值用 3 个数值表示，对应 Channel 为 3。易于理解这里假定 N=1，那么 NCHW 和 NHWC 数据格式可以很直接的这样表达：

![RGB数据格式](./.README.assets/RGB数据格式.png)

## 参考

- https://oneapi-src.github.io/oneDNN/dev_guide_understanding_memory_formats.html
- [图解 NCHW 与 NHWC 数据格式](https://blog.csdn.net/thl789/article/details/109037433)