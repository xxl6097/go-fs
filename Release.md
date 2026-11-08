### ✨ 新增功能

**公开访问目录改为运行时配置**

* 公开目录不再通过启动参数配置，改为在「设置」界面动态修改，支持运行时
  增删公开路径并落盘到 `settings.json`；
* `/fs/settings` 接口新增 `public_dirs` 字段，前端可通过此接口读取和修改；
* 已登录的受限账号也能访问公开目录，不会因为登录而失去对公开路径的访问；
* 升级 gofs 到 v0.1.15。

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
