[TOC]

# C++11 新特性

## std::function

C++11 中的 `std::function` 是一个通用的函数包装器，可以包装任何可调用的目标（函数、函数指针、成员函数指针、仿函数等），并以统一的方式进行调用。

`std::function` 的定义如下：

```c++
template<class R, class... Args>
class function<R(Args...)>;
```

其中，`R` 是返回类型，`Args` 是参数类型。可以使用 `std::function` 来定义一个函数类型，就像定义普通函数指针一样：

```c++
std::function<int(int, int)> myFunc;
```

这定义了一个函数类型 `myFunc`，接受两个 `int` 类型的参数，返回一个 `int` 类型的值。

然后，可以将任何可调用的目标（函数、函数指针、成员函数指针、仿函数等）赋值给该 `std::function` 对象：

```c++
#include <functional>

int add(int a, int b) { return a + b; }

struct MyStruct {
  int operator()(int a, int b) { return a * b; }
};

int main() {
  std::function<int(int, int)> myFunc1 = add;
  std::function<int(int, int)> myFunc2 = MyStruct();
  std::function<int(int, int)> myFunc3 = [](int a, int b) { return a - b; };
  // ...
  return 0;
}
```

在这个示例中，`myFunc1` 被赋值为一个普通函数 `add`，`myFunc2` 被赋值为一个仿函数 `MyStruct()`，`myFunc3` 被赋值为一个 Lambda 表达式。

最后，可以像调用普通函数一样使用 `std::function` 对象：

```c++
int result1 = myFunc1(1, 2); // 调用普通函数 add
int result2 = myFunc2(3, 4); // 调用仿函数 MyStruct
int result3 = myFunc3(5, 6); // 调用 Lambda 表达式
```

`std::function` 提供了一种方便的方式，可以将不同类型的可调用目标封装为相同的函数类型，并以统一的方式进行调用。
