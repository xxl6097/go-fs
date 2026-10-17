### 🐛 Bug 修复

* 修复 macOS 下升级「提示失败、其实已升级」：升级到最后拉起 restart 子进程时报 `RunChildProcess错误: fork/exec …: operation not permitted`。原因是给子进程同时设置了 `Setsid` 与 `Setpgid`，而 Go 在 macOS 上先执行 `setsid()`、再执行 `setpgid(0,0)` —— 此时子进程已经是新会话的首进程，POSIX 规定会话首进程不能改变自己的进程组，`setpgid` 返回 `EPERM`，exec 直接起不来。现只保留 `Setsid`（与 Linux 侧一致）；
* 该错误与被执行的程序本身无关：二进制其实已经替换成功，失败的是随后的「重启子进程」。macOS 上的「重启」「卸载」也走同一条路径，此前同样会失败，本次一并修复；
* FreeBSD 的同类实现存在完全相同的组合，一并修复；

---
### 🚀 github加速

```
[
  "https://ghfast.top/",
  "https://gh-proxy.com/",
  "https://ghproxy.1888866.xyz/",
  "https://github.moeyy.xyz/"
]
```
