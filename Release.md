### ✨ 新增功能

* 「已是最新版本」旁边新增「查看 GitHub Release」入口：检查更新发现已是最新时，此前只有一句话，管理员想去看看这个版本改了什么、有没有预发布版却没有入口。现在升级器若实现了可选接口 `gofs.ReleasesPage`，检查更新的响应会多带一个发布页地址，页面就在那句话旁边挂一个直达发布页的链接（新窗口打开）；
* `githubupgrade` 已实现该接口：由 `Owner` / `Repo`（以及 `APIBase`）自动推导发布页 —— github.com 用 `https://github.com/<owner>/<repo>/releases`，GitHub Enterprise 会去掉 API 根的 `/api/v3` 后再拼接；自建或镜像站推导不出来时，可用 `Options.ReleasesURL` 显式指定；
* 该扩展是可选的：没实现 `ReleasesPage` 的升级器，行为与接口字段完全不变（不下发该地址、也不显示链接）；

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
