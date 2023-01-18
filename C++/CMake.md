# CMake

生成动态库和静态库

```cmake
# 默认生成静态库
ADD_LIBRARY(hello ${LIBHELLO_SRC})
# 上面的代码可以写为
ADD_LIBRARY(hello STATIC ${LIBHELLO_SRC})
# SHARED 表示生成动态库
ADD_LIBRARY(hello SHARED ${LIBHELLO_SRC})
```

### 参考文章

- [CMake 入门实战](https://www.hahack.com/codes/cmake/)