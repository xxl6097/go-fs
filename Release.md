### 🐛 Bug 修复

* 修复分享页「页面能看，但点任何操作都跳登录页」：分享地址有两种形式 —— 查询串 `?share=<token>` 与路径 `/s/<token>/…`，而**「分享」按钮生成的正是路径形式**（为了让 wget 递归下来文件名干净）。此前只在查询串形式下发分享 cookie，路径形式打开的页面拿不到 cookie：页面本身能渲染（初始数据由服务端注入），但页面里随后的每个 XHR（列子目录、预览、上传、删除）都会 401，前端收到 401 就弹登录框。现改为「本次令牌不是从 cookie 来的」即下发，两种形式都覆盖，同时已带 cookie 的请求不再重设，避免每个 XHR 都带一次 Set-Cookie；

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
