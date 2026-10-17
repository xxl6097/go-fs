### 🐛 Bug 修复

* 修复 Windows 下服务「能安装成功、却启动不了」：`Dependencies` 此前被无条件设置成 systemd 的指令行（`After=local-fs.target network.target` / `Requires=local-fs.target`），而 Windows 会把每一项当作**依赖的服务名**交给 SCM，启动时找不到名为 `After=local-fs.target network.target` 的服务，直接报 “The dependency service does not exist or has been marked for deletion.”。现只在 Linux 上设置该系统项；
* 修复 Windows 下卸载 / 升级时删除自身二进制失败：运行中的 exe 在 Windows 上无法删除（Access is denied），现在 Windows 直接跳过这一步；Linux / macOS 仍允许 unlink 正在运行的文件，行为不变；
* 修复服务启动失败时日志出现 `%!v(MISSING)`：格式化串里留了一个 `%v` 占位却没有对应参数（错误详情已由 `zap.Error` 输出），现去掉该占位；

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
