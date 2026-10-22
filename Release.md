### ✨ 新增功能

* HTML 文件「打开」时允许页面脚本运行：CSP 由 `sandbox` 调整为 `sandbox allow-scripts` —— 带 JS 的静态页面终于能正常跑起来；同时**没有**加 `allow-same-origin`，文档仍处于不透明源，读不到本站 cookie / localStorage，也发不出带凭据的同源请求，存储型 XSS 的第二道防线不变；
* SVG / XML 维持更严格的 `sandbox`，连脚本也不开放；
* 需要脚本与同源能力都完整保留时，部署方可显式 `--no-html-sandbox`（等于主动承担存储型 XSS 风险）；

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
