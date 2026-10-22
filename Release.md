### ✨ 新增功能

* 设置页把「监听地址与端口」收进默认折叠的「高级网络设置」，摘要直接显示当前端点（`0.0.0.0:5000` 或 `unix socket`）—— 这一项低频但高风险（牵扯 Docker 端口映射、防火墙、systemd、反向代理 upstream），不再一直占着表单；
* 传输进度对读屏更友好：进行中的进度提示是 `role="progressbar"`，`aria-valuenow` / `aria-valuetext` 随进度更新（可按需读取），并关闭 `aria-live`；完成或失败时才切回 `polite`，只播报一次最终结果；

### 🐛 Bug 修复

* 修复进度提示频繁打断读屏：此前进度 toast 跟着普通 live region 每次更新都播报，读屏用户每传一块就被念一遍，基本无法使用页面；现在进行中不播报，只在完成 / 失败时播报最终状态，并同步摘掉残留的进度条语义与数值；

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
