### 🐛 Bug 修复

**升级兼容**

* 修复从旧版本原地升级后启动失败（`gofs: --auth-lock-max-delay 必须为正时长`）：旧二进制签入的 gob 配置没有新版本新增的字段，此前从零值解码导致 `AuthLockMaxDelay=0`、过不了 gofs 校验；现在解码前先铺默认值，缺失的新字段保留默认值，配置里已有的字段仍按旧值覆盖。

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
