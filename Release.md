### ✨ 新增功能

* 新增 `--allow-anonymous-admin`（环境变量 `GOFS_ALLOW_ANONYMOUS_ADMIN`，库模式 `Config.AllowAnonymousAdmin`）：未启用鉴权时**显式**放开 HTTP 管理接口。默认关闭 —— 无鉴权只读 / 上传仍是常见用法，但不该顺带把「挂载任意目录、改监听、换服务根、替换二进制、卸载服务」交给所有能连上端口的人；
* 启动 banner 会打印「管理接口：已关闭 / 已放开」，运维点开「设置」看到 403 时不必靠猜；

### 🐛 Bug 修复

**安全**

* 修复编码 `..` 导致的 ACL 越权：`ServeMux` 只对**转义后**的路径做规范化，`/docs/%2e%2e/secret.txt` 不会触发重定向，而 `r.URL.Path` 已经解码成 `/docs/../secret.txt`。此前鉴权只去掉尾斜杠，这个路径在 ACL 眼里仍属于 `/docs`，落盘时却被折叠成 `/secret.txt` —— 受限账号因此能读写自己目录之外的文件。现在 HTTP 入口统一先 `CleanURLPath`，鉴权侧用 `path.Clean` 折叠 `.` / `..` 与重复斜杠，规则本身也在归一化前就拒绝原始 `..`；读取、上传、删除与跨站点穿越都有回归测试；
* 修复 `v0.0.88` 上述修复引入的路径口径错位：鉴权层不再把 `\` 当分隔符转换 —— Unix 上它是合法的文件名字符，而落盘解析不会转换它，鉴权层转换会造成「鉴权看 `/docs/secret`、实际打开 `/docs\secret`」的错位（Windows 上落盘会把它当分隔符，鉴权只会更严格，不会放行）；
* 修复未启用鉴权时管理面完全敞开：没配账号不等于「所有访问者都是管理员」。`PUT /fs/settings`、`/fs/keys`、`/fs/users`、`/fs/upgrade`、`/fs/restart`、`/fs/uninstall` 默认一律 403，页面上也不显示「设置」入口；`GET /fs/settings` 仍可调用，但只返回上传上限等非敏感字段，不再下发服务根、可切换范围与站点列表；

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
