### ✨ 新增功能

* **目录级 TOCTOU 防护（需要 Go 1.25+）**：目录列举改为读取授权时打开的目录 fd；搜索通过 `os.Root.FS()` 遍历；建目录逐级使用 `os.Root.Mkdir`，单条删除 / 递归删除 / 重命名分别使用 `os.Root.Remove` / `RemoveAll` / `Rename`；跨设备复制逐文件走安全打开与安全建目录。即使中途把某一级目录换成根外软链，也不会在根外建目录、删文件或把整棵树搬走；
* HTTPS 根路径下的分享 Cookie 升级为按 scope 区分：根范围使用 `__Host-gofs_share`，子路径分享使用 `__Host-gofs_share_<scope 摘要>`（Secure + Path=/ + 无 Domain），选择时按令牌自身的 scope 匹配当前路径；同站子域无法再用 Domain Cookie 覆盖或伪造；
* `githubupgrade` 现在**默认要求 Ed25519 签名**：未配置 `PublicKey` / `PublicKeyText` 时 `New` 会直接报错；确实无法签名时需显式设置 `AllowUnsigned`（示例 `-allow-unsigned`）才退回只校验 SHA-256；
* `go.mod` 提升到 Go 1.25；README 同步构建要求、目录级 TOCTOU 与默认签名策略。

### 🐛 Bug 修复

**授权与目录操作**

* 修复目录列举的 TOCTOU：此前 ACL 校验后按路径重新打开目录，检查完若被换成软链，可能列出另一棵树；现在直接读取授权时打开的目录 fd；
* 修复搜索遍历可被根外软链带出服务根：递归改为通过 `os.Root.FS()` 进行；
* 修复创建目录可被中间层软链改道：`MkdirAll` 改为在挂载根 fd 下逐级 `os.Root.Mkdir`；
* 修复递归删除 / 单条删除 / 重命名可被目录替换影响：分别改为 `os.Root.RemoveAll` / `Remove` / `Rename`，批量重命名的回滚也走同一实现；
* 修复「目录是否为空」与跨设备复制时按路径二次打开：统一使用安全打开的文件 / 目录句柄；
* 修复只读打开在 Windows 上对目录语义不一致：读取路径统一使用 `os.Root.Open`。

**分享 Cookie**

* 修复 HTTPS 根路径下旧名 `gofs_share` 仍被采信：升级后根路径部署只认 `__Host-` 名字，旧 Cookie 无法再覆盖；子路径分享也迁移到 `__Host-` + scope 摘要，并校验令牌 scope 与当前路径匹配；
* 修复按 Cookie 选择分享令牌时累加使用次数并触发落盘：新增 `uploadkey.Peek`，只做校验不记统计；子路径 / 纯 HTTP 部署继续使用旧 Cookie 名字。

**升级**

* 修复无签名更新仍可默认通过：现在必须配置公钥；换掉 release 里的 `checksums.txt` 与二进制不再能因“只校验 SHA-256”而绕过；
* 新增 HTTPS 根路径旧分享 Cookie 不被采信的回归测试。

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
