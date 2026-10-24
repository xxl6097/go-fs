### ✨ 新增功能

* **明文 HTTP 默认不承载账号凭据**：启用鉴权后，Basic / 查询串口令只接受本进程 HTTPS、可信代理声明 `X-Forwarded-Proto: https`、或本机回环 / Unix socket；其余远程明文连接直接返回 403；
* 新增 `--allow-insecure-auth`（`GOFS_ALLOW_INSECURE_AUTH` / 库选项 `WithAllowInsecureAuth`）：受控内网确实需要明文时显式放开；`--cookie-secure` 可作为「外部由反代终止 TLS」的声明；
* banner 新增「传输安全」一行，说明当前部署接受凭据的连接类型（HTTPS / 可信代理 / 仅回环 / 已允许明文）；
* `/fs/health` 在启用鉴权后要求登录，不再匿名暴露版本号与存活信息；前端重启 / 卸载轮询改为「拿到小于 500 的 HTTP 响应即视为进程已回来」，重启后 401 会引导重新登录。

### 🐛 Bug 修复

**传输与会话安全**

* 修复明文 HTTP 泄露账号密码：Basic 只是 Base64，同网段抓包即可还原；现在非 HTTPS / 非可信代理 / 非本机回环的凭据请求返回 403，并提示改用 HTTPS 或显式开启 `--allow-insecure-auth`；
* 修复抓包得到的会话 Cookie 可在明文连接上复用：不安全连接不再采信会话，并顺手清掉旧 Cookie；请求仍按匿名继续，登录页可以正常打开；
* 修复反向代理部署默认无法登录：未配置 `--trusted-proxies` 时服务端无法确认外部是 HTTPS，会拒绝 Basic；README 的 nginx 说明补充该约束与 `X-Forwarded-Proto` 配置；
* 修复 `/fs/health` 匿名暴露版本与存活信息：启用鉴权后同样要求有效凭据；外部探针需带账号并通过 HTTPS / 本机回环访问；
* 新增明文 HTTP 拒绝 Basic / 查询串凭据、忽略并清除会话 Cookie、回环放行、可信代理协议判定、`--allow-insecure-auth` 逃生开关与 health 鉴权的回归测试。

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
