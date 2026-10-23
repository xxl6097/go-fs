### ✨ 新增功能

* **浏览器不再保存密码**：登录成功后身份由服务端下发的 HttpOnly 会话 Cookie 承担，密码不再写入 `sessionStorage` / `localStorage`，也不在登录后留在 JS 内存里；页面刷新 / 重新打开时用 `GET /fs/auth`（不带 Basic）换回身份、权限与配置。「保持登录」只决定 Cookie 是浏览器会话级还是 30 天，不影响密码存储；
* 设置、用户与密钥管理接口现在都接受 HttpOnly 会话身份；页面刷新后不再因为没有明文密码而无法管理密钥 / 分享；
* HTTPS 根路径下会话 Cookie 优先使用 `__Host-gofs_sid`，分享 Cookie 使用 `__Host-gofs_share`（Secure + Path=/ + 无 Domain），同站子域无法用 Domain Cookie 覆盖或伪造；从旧版本升级到 HTTPS 根路径部署需重新登录一次；
* 新增 TOCTOU 安全打开：授权时记录目标的 inode（新建时记录最近已存在祖先），打开时用 `os.Root` 从挂载根逐段解析并用 `SameFile` 复核，新建文件使用 `O_EXCL`；检查后路径被换成软链或另一个文件时返回 409，而不是按旧权限读写；
* `githubupgrade` 支持校验和文件的 Ed25519 分离签名：配置 `PublicKey` / `PublicKeyText`（hex 或 base64）后，必须存在 `checksums.txt.sig` 且验证通过，否则拒绝升级；`examples/upgrade` 新增 `-public-key`；
* 新增 `api.md`：整理 `PUT` / `POST` 上传、`/up` 一键脚本、`/fs/keys` 密钥接口与认证方式的完整示例。

### 🐛 Bug 修复

**认证与会话**

* 修复明文密码可被 XSS 读取：此前 Basic 密码保存在 `sessionStorage` / `localStorage` 并在每个请求里由 JS 附加；现在只保留用户名标记，旧版本残留的密码会在启动时主动清除；
* 修复 cookie tossing：同站子域此前可以用 Domain Cookie 覆盖会话 / 分享 Cookie；HTTPS 根路径现在只认 `__Host-` 前缀的 Cookie，旧名字不再被采信；
* 修复共享 / 会话 Cookie 清理不完整：登出与路径前缀迁移会同时清理新旧两种名字，避免升级后残留一张仍可用的旧会话票；
* 修复设置 / 密钥接口依赖 Basic 密码：页面刷新后只有会话 Cookie 时，设置、用户、密钥与分享创建仍可正常工作。

**授权与文件打开**

* 修复 ACL 检查与真正打开之间的 TOCTOU：检查完成后路径若被换成软链或另一个 inode，此前会按旧权限打开目标；现在下载、Office 预览、文本读写、表单上传、压缩包读取 / 解压、zip 打包与 index.html 都使用授权快照 + `os.Root` 打开，不一致返回 409；
* 修复 multipart 上传的 rename 快路径可被目录软链改道：现在改用授权时解析出的真实路径作为 rename 目标，并复核父目录快照；
* 修复压缩包按路径二次打开：`archive.OpenFile` 复用已通过快照校验的文件句柄，避免检查完再按路径重开时被替换；
* 修复新建目标可被抢先替换为软链：新建写入统一使用 `O_EXCL`，授权后路径上出现软链 / 已存在文件会直接失败。

**解压与升级**

* 修复 zip 压缩比可被声明值稀释：改用中央目录数据偏移算出的**物理跨度**计算压缩比，不再相信可伪造的 `CompressedSize64`；
* 修复升级校验和仍可能被同一 release 篡改：配置公钥后强制验证 `checksums.txt.sig` 的 Ed25519 签名；签名或校验和响应体也有大小上限。

**其它**

* 新增会话刷新、`__Host-` Cookie、TOCTOU 替换检测、软链新建拒绝与签名校验等回归测试。

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
