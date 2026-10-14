### 🐛 Bug 修复

* 修复分享页里点开文件、预览、解压时弹登录框：这几个动作走的是 `/fs/text`、`/fs/office`、`/fs/extract` 等 `/fs/*` 接口，而它们**不在分享路径（如 `/docs/`）之下** —— 分享 cookie 的 Path 正好限定在分享范围内，浏览器根本不会把 cookie 带到 `/fs/*`，于是请求 401、前端收到 401 就弹登录框。现给这些接口调用统一补上 `?share=<token>`；
* 修复分享页进入子目录后地址栏丢掉分享前缀：页面内部一律按 `state.path`（`/docs/…`）拼 URL，一进子目录地址栏就退化成 `/docs/sub/` —— 把地址复制给别人打开只会看到登录框（那条路径不属于任何分享）。现文件空间的地址一律带上 `/s/<token>` 前缀，复制出去仍是能直接打开的分享地址，下载、预览、`<img>`/`<video>` 这类不走接口封装的直接地址也因此带上了令牌；页面内解析当前目录时会先剥掉这层前缀，再按 `state.path` 处理；

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
