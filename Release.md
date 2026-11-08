### 🐛 Bug 修复

**修复设置面板变量作用域导致崩溃**

* 变量 `publicDirsData` 原本在 `fill()` 内用 `var` 声明，
  但在外层的 `dirty()` 和 `save()` 中引用，导致 `ReferenceError`；
* 表现为**点击保存按钮无响应**、**无法切换 tab**（密钥、用户与权限等），
  整个设置面板崩溃；
* 修复：`publicDirsData` 移至 `buildGeneralTab` 外层作用域声明，
  `fill()` 内改为直接赋值；导出按钮的 `copyText` 改为 `G.copyText`
  （跨文件引用）；
* 升级 gofs 到 v0.1.21。

---
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
