### ✨ 新增功能

* **明文 HTTP 认证默认开启，HTTP 部署开箱即用**：账号密码与会话可在明文连接上直接使用，不再需要额外开关；
* 需要更严格时使用 `--allow-insecure-auth=false`（或 `GOFS_ALLOW_INSECURE_AUTH=false`、库选项 `WithAllowInsecureAuth(false)`）：此时只接受本进程 HTTPS、可信代理声明的 `X-Forwarded-Proto: https`、本机回环或 Unix socket；
* 收紧模式下，其余明文连接上的 Basic / 查询串口令返回 403，会话 Cookie 不采信并清除（请求按匿名继续，登录页仍能打开）；
* banner 的「传输安全」行调整判断顺序：先报 HTTPS / `--cookie-secure` / 可信代理，最后才提示默认允许明文认证。

### 🐛 Bug 修复

**配置与传输**

* 修复默认策略导致 HTTP 部署无法登录：现在默认允许明文认证，保持开箱即用；需要强制 HTTPS 的部署显式收紧；
* 修复 HTTPS 部署的 banner 被误报为「允许明文认证」：先判断 TLS / 反代声明，再回落到默认提示；
* 修复 README 与 nginx 说明把「未配置可信代理」写成无条件 403：该后果只在 `--allow-insecure-auth=false` 收紧后发生；
* `api.md` 的错误码说明同步为：明文提交账号密码仅在收紧模式下返回 403；
* 新增默认配置允许远程明文 Basic 的回归测试；严格模式仍由传输安全测试覆盖。

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
