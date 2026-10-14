# pkg/github 使用说明

从 GitHub Releases 读取最新版本、和本地版本比对、解析代理地址、给出下载链接。
只依赖 `net/http`，没有引入 go-github 之类的 SDK，所有接口都是直接打 REST API。

---

## 一、前置约定

`CheckUpgrade` 依赖三条写死的约定，Release 不合规就拿不到结果。

### 1. 资源（Asset）命名必须带版本号

本地文件名里要有 `vX.Y.Z`，发布的新资源名必须是「把本地文件名里的版本号换成新 tag」：

| 项 | 值 |
|---|---|
| 本地文件名 `fullName` | `aatest_v0.0.20_darwin_arm64` |
| Release tag | `v0.0.21` |
| 计算出的目标资源名 | `aatest_v0.0.21_darwin_arm64` |

实现见 `utils.ReplaceNewVersionBinName`（正则 `_v\d+\.\d+\.\d+_` → `_<tag>_`）。
`fullName` 通常直接用 `version.BinName`，由 `build.sh` 注入，命名规则是 `${appname}_${version}_${os}_${arch}`。

匹配是**忽略大小写的全等比较**，所以资源名对不上（比如多了 `-debug` 后缀）就会 `fullUrl` 为空。

### 2. 差量包命名为 `<全量包名>.patch`

即上例中的 `aatest_v0.0.21_darwin_arm64.patch`，和全量包放在同一个 Release 里。

### 3. Release 正文里用代码块写代理前缀列表

`Request()` 会调用 `utils.ParseMarkdownCodeToStringArray(result.Body)`，它取正文里**第一个围栏代码块**，按 JSON 反序列化成 `[]string`：

````markdown
本次更新内容……

```json
["https://ghproxy.net/", "https://mirror.ghproxy.com/"]
```
````

解析失败或没有代码块时 `proxies` 为空，不会有报错，只是退化成直连。正文中**第一个 `---` 之前的文字**会被当作 release note 返回。

---

## 二、环境变量

| 变量 | 作用 | 说明 |
|---|---|---|
| `GITHUB_API_PROXY` | 给 **API 元数据请求**加前缀 | `https://api.github.com/repos/...` 会变成 `<proxy>/https://api.github.com/repos/...`；结尾缺 `/` 会自动补 |
| `GITHUB_CLIENT_ID` | API Basic Auth 用户名 | 由 `githubKey.dat` 注入，见下节 |
| `GITHUB_CLIENT_SECRET` | API Basic Auth 密码 | 同上 |

注意 `GITHUB_API_PROXY` 只影响 **`/releases/latest` 这一次请求**。资源下载走的是 `BrowserDownloadUrl`，由 `GetProxyUrls` 在**返回值层面**逐个拼代理前缀，两条链路互不相干。

`GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` 用来抬高 API 限流额度（未认证 60 次/小时/IP）。两者都为空时**不加** Authorization 头。

---

## 三、凭据文件 `githubKey.dat`

`LoadGithubKey()` 从下面这个路径读 gob 文件：

```
zutil.AppHome("obj")/githubKey.dat
= <TempDir>/aalog/<AppName>/obj/githubKey.dat
```

- macOS/Linux：`/var/folders/.../T/aalog/aatest/obj/githubKey.dat`、`/tmp/aalog/aatest/obj/githubKey.dat`
- Windows：`%ProgramData%\aalog\<AppName>\obj\githubKey.dat`

结构体是 `model.GithubKey{ ClientId, ClientSecret }`。

**本仓库只读不写**——`LoadGithubKey` 是唯一的消费者，全仓库没有任何地方生成这个文件。要启用认证，得在你的程序里自己落盘：

```go
import (
    "path/filepath"
    "github.com/xxl6097/glog/pkg/zutil"
    "github.com/xxl6097/go-service/pkg/github/model"
    "github.com/xxl6097/go-service/pkg/utils"
)

fpath := filepath.Join(zutil.AppHome("obj"), "githubKey.dat")
_ = utils.SaveWithGob(model.GithubKey{
    ClientId:     "Iv1.xxxxxxxxxxxx",
    ClientSecret: "yyyyyyyyyyyyyyyyyyyy",
}, fpath)
```

两个字段必须**都非空**才会生效，否则会把环境变量显式置空（即清掉之前设过的值）。

> ⚠️ 路径落在 `os.TempDir()` 下。macOS 会定期清理 `/var/folders`，Linux 重启清 `/tmp`，所以这份凭据可能莫名其妙消失。长期部署建议直接在进程环境里设 `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET`（`LoadGithubKey` 在读到空文件时会把它们覆盖成空串，所以两者不要混用）。

---

## 四、使用步骤

### Step 1 — 取单例、指定仓库

```go
api := github.Api().SetName("xxl6097", "go-service")
```

- `Api()` 返回进程内单例，**第一次调用时**执行 `LoadGithubKey()`（`sync.Once`，之后不再重读环境变量）。
- `SetName` 空字符串会被忽略，支持链式调用。
- **必须在任何会走网络的方法之前调用**，否则 `defaultRequest` 直接返回 `请指定github的用户名` / `请指定github的仓库名`。

### Step 2 — 检查是否有新版本

```go
data, err := api.CheckUpgrade(version.BinName)
if err != nil {
    // 已是最新版时 err 为 "已是最新版本～"
    log.Println(err)
    return
}
// 有新版：
fullUrl  := data["fullUrl"].(string)      // 全量包下载地址，可能为空
patchUrl := data["patchUrl"].(string)     // 差量包下载地址，可能为空
notes    := data["releaseNotes"].(string) // markdown 文本，可直接渲染
```

`CheckUpgrade` 内部做的事：

1. `request()` → `GET /repos/{user}/{repo}/releases/latest`（satisfy `GITHUB_API_PROXY`、Basic Auth）
2. 用 `utils.GetVersionByFileName(fullName)` 从本地文件名里正则抠出 `vX.Y.Z`
3. `utils.CompareVersions(最新tag, 本地版本)`，返回 `>0` 才算有新版
4. 在 `Assets` 里按名字找全量包和 `.patch`，顺便把大小格式化成 `1.23 MB`
5. 拼出 `releaseNotes`

### Step 3 —（可选）拿代理 URL 再下载

`GetProxyUrls` 只是**拼字符串**，不发起请求：

```go
urls := api.GetProxyUrls(fullUrl)
// ["https://ghproxy.net/https://github.com/.../aatest_v0.0.21_darwin_arm64",
//  "https://mirror.ghproxy.com/https://github.com/...",
//  "https://github.com/..."]          ← 原始地址永远排在最后
filePath := utils.DownloadFileWithCancelByUrls(urls)  // 并发竞速，返回 "" 表示失败
```

> `DownloadFileWithCancelByUrls` 走的是 `utils.DynamicSelect`：给每个 URL 起一个 goroutine
> **同时**下载，`reflect.Select` 收第一个落到 channel 上的结果就 `cancel()` 并返回。
> 所以它是「谁先返回用谁」而不是按顺序重试——某个镜像快速失败也可能抢先返回空串。

`GetProxyUrls` 的返回值还依赖 `proxies` 是否已被填充，而 `proxies` 只在 `Request()`
（即 `CheckUpgrade` / `GetModel` / `GetDownloadUrl(s)` / `defaultRequest`）里赋值。
**所以要么先调一次 `CheckUpgrade`，要么接受它退化成只返回原始地址。** 单例状态下第一次调用通常就是这样。

### 或者：不用版本比对，直接按条件找资源

```go
// 拿最新 release 里第一个满足条件的资源地址
url := api.GetDownloadUrl(func(tag string, a *model.Assets) bool {
    return strings.Contains(a.Name, "darwin_arm64")
})

// 拿全部满足条件的
urls := api.GetDownloadUrls(func(tag string, a *model.Assets) bool {
    return strings.HasSuffix(a.Name, ".patch")
})
```

`fn` 传 `nil` 时不会匹配到任何资源（内部有 `fn != nil` 判断），返回空。

---

## 五、API 速查

| 方法 | 网络请求 | 返回 |
|---|---|---|
| `Api() *githubApi` | 首次调用无 | 单例，内部 `LoadGithubKey()` |
| `SetName(user, repo) *githubApi` | 无 | 自身，便于链式 |
| `Request(user, repo) (*model.GitHubModel, error)` | 有 | 完整 release 结构；副作用：填充 `proxies` |
| `CheckUpgrade(fullName) (map[string]any, error)` | 有 | `fullUrl` / `patchUrl` / `releaseNotes` |
| `GetModel() *model.GitHubModel` | 有 | 失败返回 `nil`（已打日志） |
| `GetDownloadUrl(fn) string` | 有 | 第一个命中的 `BrowserDownloadUrl` |
| `GetDownloadUrls(fn) []string` | 有 | 全部命中的 `BrowserDownloadUrl` |
| `GetProxyUrls(fileUrl) []string` | 无 | 代理前缀 × 列表 + 原始地址 |
| `LoadGithubKey()` | 无（读文件） | 设置/清空两个环境变量 |

失败时的日志级别：`request` 用 `Warn`，`GetModel` / `GetDownloadUrl(s)` 内部 `Warn` 后返回零值——**注意这三个方法把错误吞掉了**，调用方拿不到失败原因，只能看到空结果。

---

## 六、完整示例

```go
package main

import (
    "fmt"

    "github.com/xxl6097/go-service/pkg/github"
    "github.com/xxl6097/go-service/pkg/utils"
    "github.com/xxl6097/go-service/pkg/version"
)

func main() {
    api := github.Api().SetName("xxl6097", "go-service")

    data, err := api.CheckUpgrade(version.BinName)
    if err != nil {
        fmt.Println("无需升级:", err)
        return
    }
    fmt.Println(data["releaseNotes"])

    // 优先差量包，没有再退回全量包
    url, _ := data["patchUrl"].(string)
    if url == "" {
        url, _ = data["fullUrl"].(string)
    }
    if url == "" {
        fmt.Println("Release 里没找到对应资源")
        return
    }

    filePath := utils.DownloadFileWithCancelByUrls(api.GetProxyUrls(url))
    if filePath == "" {
        fmt.Println("下载失败") // 所有镜像并发竞速，谁先返回用谁
        return
    }
    fmt.Println("已下载到:", filePath)
}
```

---

## 七、已知问题 / 注意点

1. **差量升级的「只差一个版本」判断没有生效。**
   `CheckUpgrade` 里写着
   ```go
   if hasNewVersion != 1 { patchUrl = "" }
   ```
   注释说「版本之间只有相差一个版本号才能差量升级」，但 `utils.CompareVersions` 只会返回
   `1 / -1 / 0`（遇到第一个不同段就直接 return），所以 `hasNewVersion > 0` 必然等于 `1`，
   这个分支永远不会命中，`patchUrl` 也就永远不会被清空。

2. **`releases/latest` 取不到预发布版本。** GitHub 的这个端点会跳过 `prerelease` 和 `draft`。
   要覆盖预发布得自己改成 `GET /releases` 再挑第一个。

3. **单例 + `sync.Once` 的组合不太好测。** `LoadGithubKey` 只在首次 `Api()` 时跑一次，
   之后改环境变量不会再读；`userName` / `repoName` / `proxies` 都是实例状态，
   多仓库场景要反复 `SetName`，且 `proxies` 会被下一次 `Request` 覆盖。

4. **`GetModel` / `GetDownloadUrl` / `GetDownloadUrls` 吞错误**，只打日志、返回零值，
   调用方无法区分「没有匹配资源」和「请求失败」。

5. **`GetProxyUrls` 有隐式调用顺序要求**，必须先有成功过的 `Request`，否则代理列表为空（见 Step 3）。

6. `ParseMarkdownCodeToStringArray` 只认**第一个**代码块，且要求是合法的 JSON 字符串数组；
   Release 正文里如果有别的代码块写在前面（比如贴了段 shell），代理列表就会解析失败并静默降级。
