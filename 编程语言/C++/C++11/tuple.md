# tuple

C++11 标准新引入了一种类模板，命名为 tuple（中文可直译为元组）。tuple 最大的特点是：实例化的对象可以存储任意数量、任意类型的数据。

tuple 的应用场景很广泛，例如当需要存储多个不同类型的元素时，可以使用 tuple；当函数需要返回多个数据时，可以将这些数据存储在 tuple 中，函数只需返回一个 tuple 对象即可。

```c++
#include <iostream>
#include <tuple>

int main() {
  int size;
  // 创建一个 tuple 对象存储 10 和 'x'
  std::tuple<int, char> mytuple(10, 'x');
  // 计算 mytuple 存储元素的个数
  size = std::tuple_size<decltype(mytuple)>::value;  // size = 2
  // 输出 mytuple 中存储的元素
  // 输出 10 x
  std::cout << std::get<0>(mytuple) << " " << std::get<1>(mytuple) << std::endl;
  // 修改指定的元素
  std::get<0>(mytuple) = 100;
  std::cout << std::get<0>(mytuple) << std::endl;  // 输出 100

  // 使用 makde_tuple() 创建一个 tuple 对象
  auto bar = std::make_tuple("test", 3.1, 14);
  // 拆解 bar 对象，分别赋值给 mystr、mydou、myint
  const char *mystr = nullptr;
  double mydou;
  int myint;
  // 使用 tie() 时，如果不想接受某个元素的值，实参可以用 std::ignore 代替
  // 使用 tie 把 tuple 中的元素赋值到定义的变量中
  std::tie(mystr, mydou, myint) = bar;
  // 输出 test 3.1 14
  std::cout << mystr << " " << mydou << " " << myint << std::endl;
  // std::tie(std::ignore, std::ignore, myint) = bar;  //只接收第 3 个整形值

  // 将 mytuple 和 bar 中的元素整合到 1 个 tuple 对象中
  auto mycat = std::tuple_cat(mytuple, bar);
  size = std::tuple_size<decltype(mycat)>::value;
  std::cout << size << std::endl;  // 输出 5
  return 0;
}

```

## 参考文章

- [C++11 tuple 元组详解](http://c.biancheng.net/view/8600.html)