### 🐛 Bug 修复

**修复 Mermaid 宽图被压缩、显示不完整**

* 修复 mermaid 默认给 SVG 输出 `width=100%` + `max-width`，导致宽图被压到
  预览栏宽度、文字缩成蚂蚁看不清的问题：改为渲染后按 `viewBox` 恢复自然尺寸，
  宽图交给卡片横向滚动；
* 修复用 `text-align: center` 居中 SVG 时，溢出图左侧会被裁掉的问题：
  改用 `margin: 0 auto`（窄图居中、宽图左对齐后横向滚动）；
* 无头 Chrome 自检：宽图 / 窄图 / 高图 / 时序图均保持自然尺寸、无纵横裁切；
* 升级 gofs 到 v0.1.13。

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
