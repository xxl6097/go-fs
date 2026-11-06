### ✨ 新增功能

**Markdown 预览支持 Mermaid 图**

* ` ```mermaid ` 围栏代码块在预览里渲染成 SVG（流程图 / 时序图 / 甘特图 /
  饼图 / 状态图等），内嵌官方 mermaid v11（MIT），不依赖 CDN、离线可用；
* 只在第一次遇到 mermaid 块时才懒加载渲染库，普通 Markdown 不额外付出体积；
* 渲染走 mermaid `securityLevel: strict`，SVG 输出先消毒再插入预览，
  图文本仍按普通代码块转义，渲染失败时回退显示源码并标出错误；
* 图随浅色 / 深色主题自动切换 `default` / `dark` 主题；
* 升级 gofs 到 v0.1.11。

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
