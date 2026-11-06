### 🐛 Bug 修复

**优化 Mermaid 预览布局与配色**

* 修复 Mermaid 图按整页宽度排布、再塞进预览栏，导致节点文字被压缩挤成一团
  的问题：改为把预览里的宿主元素作为容器传给 mermaid，按预览栏实际宽度排布；
* 修复渲染时容器尚未挂进 DOM、宽度为 0 导致布局错乱：改为等预览挂载后再渲染；
* 弃用 mermaid 默认的紫 / 粉彩色主题，改用 `neutral` + 对齐 gofs 配色的
  `themeVariables`（浅色白底 / 深色深底 + 主题蓝描边），随主题切换文字与线条色；
* SVG 卡片居中、宽图横向可滚，渲染前保留源码块作回退；
* 升级 gofs 到 v0.1.12。

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
