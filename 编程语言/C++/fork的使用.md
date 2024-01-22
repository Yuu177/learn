[TOC]

# fork 的使用

## Description

1. 基本的 fork 使用
2. 父进程如何守护子进程
3. 父进程意外退出后，子进程如何避免成为孤儿进程
4. 信号处理函数 signal 的使用
5. 子进程 crash 后，如何「优雅」地退出重启，避免产生 core 文件

## 代码示例

```c++
#include <cstdlib>
#include <unistd.h>
#include <wait.h>
#include <iostream>
#include <sys/prctl.h>

constexpr auto SIG_RESTART = SIGUSR1;  // 使用自定义信号 SIGUSR1 做为子进程重启信号

void RestartBySignal(int sig) {
  std::cout << "Recv sig and raise: " << sig << std::endl;

  raise(SIG_RESTART);  // 自己给自己发送 SIG_RESTART 信号
}

constexpr auto EXIT_STATE_RESTART = 1; // 定义退出重启的状态码

void RestartByExit(int sig) {
  std::cout << "Recv sig and exit: " << sig << std::endl;

  exit(EXIT_STATE_RESTART); // 进程退出返回 EXIT_STATE_RESTART 状态码
}

// 模拟程序运行
void Run() {
  std::cout << getpid() << " fork from " << getppid() << " and run" << std::endl;

  // signal 设置一个函数来处理对应信号
  // signal(SIGABRT, RestartBySignal); // 捕获 SIGABRT 并转换为 SIG_RESTART
  signal(SIGABRT, RestartByExit);   // 捕获 SIGABRT 并 exit

  // prctl(PR_SET_PDEATHSIG, SIGKILL); // 防止子进程变成孤儿进程方法一：Only Linux，父进程退出时，会收到 SIGKILL 信号

  while (true) {
    // 防止子进程变成孤儿进程方法二
    // https://stackoverflow.com/questions/284325/how-to-make-child-process-die-after-parent-exits
    static auto ppid = getppid();
    if (getppid() != ppid) {
      std::cout << getppid() << " diff " << ppid << ", the parent process has already died" << std::endl;
      exit(0);
    }

    sleep(3);
    abort(); // Abort execution and generate a core-dump.
  }
}

void RunParentProcess(int child_pid) {
  int child_status;
  // waitpid 阻塞直到子进程退出
  waitpid(child_pid, &child_status, 0);

  // WIFEXITED(status) 用于检查子进程是否正常退出。当子进程正常退出时（通过调用 exit() 或 _exit()），
  // WIFEXITED(status) 将返回非零值，表示子进程以正常方式终止。
  // 可以使用 WEXITSTATUS(status) 来获取子进程的退出状态码，这个状态码是子进程传递给 exit() 函数的值。
  if (WIFEXITED(child_status)) {
    if (WEXITSTATUS(child_status) == EXIT_STATE_RESTART) {
      std::cout << child_pid << " recv restart exited status" << std::endl;
      goto restart;
    }
    std::cout << child_pid  << " exited with status: " << WEXITSTATUS(child_status) << std::endl;
  }

  // WIFSIGNALED(status) 用于检查子进程是否因为收到一个信号而终止。
  // 如果子进程因为收到一个信号而非正常退出（例如被 kill 发送了一个信号），那么 WIFSIGNALED(status) 将返回非零值。
  // 可以使用 WTERMSIG(status) 来获取导致子进程终止的信号编号。
  if (WIFSIGNALED(child_status)) {
    if (WTERMSIG(child_status) == SIG_RESTART) {
      std::cout << child_pid << " recv restart signal" << std::endl;
      goto restart;
    }
    std::cout << child_pid  << " terminated by signal: " << WTERMSIG(child_status) << std::endl;
  }

  exit(0); // 这里子进程意外退出父进程也一起退出，可以根据实际情况注释掉，让父进程来守护子进程

restart:
  std::cout << "Restart child process\n" << std::endl;
  sleep(1);
}

int main() {
  int child_pid = -1;
  while (true) {
    child_pid = fork();
    if (child_pid == -1) {
      std::cout << "Failed to fork." << std::endl;
      exit(-1);
    } else if (child_pid == 0) {
      // 子进程
      Run();
      return 0;
    } else {
      // 父进程
      RunParentProcess(child_pid);
    }
  }

  return 0;
}

```

## 其他

### exit 和 kill 区别

- **exit()**：
  - `exit()` 是一个库函数，用于正常终止进程的执行。
  - 调用 `exit()` 函数会使程序按照正常流程结束执行，清理资源并返回一个退出状态码。
  - 它并不是发送信号给进程，而是在程序内部执行，通常在程序已经完成其任务或者发生错误时使用。
  - 可以通过传递一个退出状态码来表示程序的退出状态，这个状态码可以被其他进程或系统检索到，例如父进程可以通过 `wait()` 系统调用获取子进程的退出状态码。
- **kill()**：
  - `kill()` 是一个系统调用，用于向指定进程发送信号。
  - 通过 `kill()` 函数可以发送各种不同的信号给指定的进程，比如终止信号 (`SIGTERM`)、强制终止信号 (`SIGKILL`)、用户自定义信号等。
  - 它用于与其他进程进行通信，通常用于控制进程的行为，例如关闭、重启、中止或触发特定操作等。
  - 调用 `kill()` 函数并不会导致进程立即终止，而是向目标进程发送信号，具体如何响应信号取决于接收进程对该信号的处理设置。

总的来说，`exit()` 是用于程序内部正常退出的函数，而 `kill()` 则是用于向指定进程发送信号，以影响其行为或状态。

### SIGABRT

进程在释放内存时出现了重复释放、访问无效内存或执行非法指令等问题，操作系统内核可能会检测到这种行为并向进程发送信号 6（SIGABRT）。这个信号通常表示出现了严重的问题，如非法指令或内存错误，操作系统会采取这种方式通知进程发生了异常情况，进而终止该进程的执行或进行相应的处理。
