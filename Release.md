### ✨ 新增功能

* 源码格式化：JSON / XML / YAML / TOML / INI / conf / .config / .env 的编辑器工具栏新增「格式化」按钮，直接改写正文（会进入未保存状态，仍可撤销 / 保存）；
* 各格式的具体行为：JSON 用 `JSON.parse` + 2 空格缩进，非法 JSON 只报错不改内容；XML 经 `DOMParser` 校验后按层级缩进，混合内容保留原内部形态，`xml:space="preserve"` 的子树不重排；YAML 保守整理块式配置的缩进（统一 2 空格），`|` / `>` 块标量的内容与行尾空白原样保留、只随块头整体平移，遇到跨行引号 / Tab 缩进 / 非块式多行标量时**拒绝格式化**而不是冒险改坏；TOML 统一 section 头与 `key = value` 空白，多行 `"""` / `'''` 原样保留；INI / conf / config 统一 section 头与键值空白，缩进续行原样保留；`.env` 统一成 `KEY=value`，保留 `export` 前缀、引号与行内注释；
* 格式化不改变原文件的换行风格（CRLF / LF）与 BOM，保存时仍按原样写回；
* `.config` 扩展名（Windows 上常见的配置文件）现在也按 INI 处理，可在编辑器里打开并格式化；

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
