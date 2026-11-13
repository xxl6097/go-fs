### ✨ 新增功能

**支持自定义访问路径前缀**

* 安装/初始化配置时新增「访问路径前缀」输入项，例如填 `/files` 后管理后台挂在 `http://host:port/files/` 下访问；
* 留空则不设置前缀，保持原行为（服务挂在根 `/`）；
* 前缀由 gofs 在启动时统一归一化（自动补前导斜杠、去尾斜杠），鉴权前会剥掉前缀，无需调整鉴权规则。

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
