### ✨ 新增功能

**WebDAV 服务**

* 新增 `/dav/` 端点，Windows「映射网络驱动器」、macOS Finder（⌘K）、
  Cyberduck、davfs2 等客户端可直接挂载；
* WebDAV 账号与网页登录账号彼此独立，每个账号单独授权可访问的目录，
  未授权路径一律拒绝；密码只存 PBKDF2 摘要；
* 「设置 → WebDAV 账号」支持新建账号、修改目录权限、重置密码、删除，
  并提供 `/api/webdav/accounts` 管理接口；
* 浏览器打开 `/dav/<目录>/` 时直接给出只读文件列表，可逐级进入目录、
  点文件下载或就地预览。

### 🐛 Bug 修复

* 修复 WebDAV 账号面板「新建账号」等按钮点击无反应：按钮字段与弹窗 API
  不匹配（`action`/`G.overlay` 并不存在），改用 `onClick` 与
  `G.openModal`；顺带去掉重复的「关闭」按钮；
* 修复只注册了账号管理接口、没有挂载 `/dav/`，导致账号建好却无法连接的问题；
* 修复挂载前缀下目录权限判断错误：请求路径 `/dav/docs` 与授权目录 `/docs`
  对不上，现在先剥离前缀再比对；
* 修复连到未授权根目录时只返回一句 403，现在响应里会列出该账号实际被授权的目录；
* 升级 gofs 到 v0.1.6。

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
