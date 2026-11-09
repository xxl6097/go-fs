### 🐛 Bug 修复

**修复 systemd 服务启动时 Huge Pages panic**

* 在 systemd 服务配置中启用 GODEBUG 环境变量，设置 `GODEBUG=madvdontneed=1`；
* 解决 Go 运行时在初始化阶段调用 `getHugePageSize()` 时的 panic 问题；
* 环境变量现在在程序启动前由 systemd 设置，确保在 Go 运行时初始化前生效；
* 移除冗余的 `os.Setenv()` 调用，使用 kardianos/service 的官方机制。

**完善设置面板公开目录反馈与数据同步**

* 保存公开目录时在 toast 中显示具体条数（「公开目录 N 条」）；
* 确保本地数据与服务端返回值同步，防止显示不一致；
* 优化用户反馈体验，让用户清楚地看到保存结果；
* 升级 gofs 到 v0.1.22。

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
