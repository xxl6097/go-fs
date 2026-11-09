### 🐛 Bug 修复

**修复安装后服务启动 panic（fatal error: index out of range / panic before malloc heap initialized）**

真正的成因不是 Huge Pages，而是**签名时把二进制撑长了**：

* 编译产物内嵌一段占位区（`cmd/app/buffer/buffer`，`//go:embed`，原为 1024 字节 `+`）；
  安装/签名时 `ukey.GenConfig` 生成 `rawKey(39B) + 4B 长度 + 配置密文`，原地替换这段占位符；
* `GenConfig` 的容量校验只比较了 `body` 与 `len(buffer)-4`，**漏算 `rawKey+4` 的 43 字节头部**，
  于是当 `data` 落在「比占位区长一点点」的区间时既不报错、也不补齐；
* 本次配置密文 998B，`39+4+998=1041 > 1024`，`GenerateBin` 把 1024 字节占位符换成 1041 字节，
  文件从 13893794 变 13893811（**+17**），但 ELF 的段表/节表偏移没有更新；
* 插入点落在最后一个 `PT_LOAD` 段内，其后所有数据整体后移 17 字节，
  运行时初始化 `osinit → getHugePageSize()` 读到错位的全局 `sysTHPSizePath`（len 变 0），
  `&sysTHPSizePath[0]` 触发 bounds panic —— 这就是那句 `getHugePageSize` 报错的来历，
  它与 Huge Pages、GODEBUG 均无关。

修复：

* 把占位区从 1024 字节扩到 **8192 字节**（`cmd/app/buffer/buffer`），给配置留足余量；
  签名数据与占位区等长，文件长度不变，ELF 不再错位；
* 移除此前误判加的 `GODEBUG=madvdontneed=1`（Go 1.16 起该开关已无意义，且改不了这个问题）；
* 验证方式：`go version`/`sha256sum` 对比、`readelf`/`debug/elf` 检查 `PT_LOAD` 是否与编译产物一致。

**完善设置面板公开目录反馈与数据同步**

* 保存公开目录时在 toast 中显示具体条数（「公开目录 N 条」）；
* 确保本地数据与服务端返回值同步，防止显示不一致；
* 优化用户反馈体验，让用户清楚地看到保存结果；
* 升级 gofs 到 v0.1.22。

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
