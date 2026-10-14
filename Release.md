### ✨ 新增功能

* 新增「卸载」能力：页面「升级、重启与卸载」里可卸载本服务 —— 停止服务并清理配置目录（密钥、用户表、运行期设置）；卸载是独立扩展点（`WithUninstaller`），与重启各管各，只重启不卸载、只卸载不重启都能各自成立；
* 服务升级 / 重启 / 卸载的适配器统一收敛到 `internal`（`Upgrader` / `Restarter` / `Uninstaller`），启动流程只接收一个服务实例，注入哪些能力由 `internal.RunCfg` 一处决定；

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
