### 🐛 Bug 修复

**公开目录绿色徽标在页面加载时即可显示**

* 之前 `public_dirs` 只随 `/fs/auth` 登录响应下发，初始页面渲染时
  `DATA.public_dirs` 为空，「公开」绿色徽标和行菜单状态要等用户触发一次
  认证后才生效；
* 修复：`pageData` 与 `restrictedPageData`（匿名/分享访问）均新增
  `PublicDirs` 字段，服务端渲染时即携带公开目录列表；
* 现在打开页面即可在已配置的公开目录旁看到绿色「公开」徽标，无需先登录；
* 升级 gofs 到 v0.1.20。

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
