package github

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-service/pkg/github/model"
	"github.com/xxl6097/go-service/pkg/utils"
	"go.uber.org/zap"

	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func LoadGithubKey() {
	fpath := filepath.Join(zutil.AppHome("obj"), "githubKey.dat")
	obj, err := utils.LoadWithGob[model.GithubKey](fpath)
	if err == nil && obj.ClientId != "" && obj.ClientSecret != "" {
		os.Setenv("GITHUB_CLIENT_ID", obj.ClientId)
		os.Setenv("GITHUB_CLIENT_SECRET", obj.ClientSecret)
	} else {
		os.Setenv("GITHUB_CLIENT_ID", "")
		os.Setenv("GITHUB_CLIENT_SECRET", "")
	}
}

var (
	instance *githubApi
	once     sync.Once
)

type githubApi struct {
	proxies            []string
	userName, repoName string
}

// Api 返回单例实例
func Api() *githubApi {
	once.Do(func() {
		instance = &githubApi{}
		LoadGithubKey()
	})
	return instance
}

func request(githubUser, repoName string) ([]byte, error) {
	baseUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", githubUser, repoName)
	proxy := os.Getenv("GITHUB_API_PROXY")
	if proxy != "" {
		if !strings.HasSuffix(proxy, "/") {
			proxy = fmt.Sprintf("%s/", proxy)
		}
		baseUrl = fmt.Sprintf("%s%s", proxy, baseUrl)
	}
	z.L().Info("request", zap.String("url", baseUrl))
	client := &http.Client{}
	req, _ := http.NewRequest("GET", baseUrl, nil)
	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if clientId != "" || clientSecret != "" {
		req.SetBasicAuth(clientId, clientSecret) // 自动 Base64 编码
	}
	resp, err := client.Do(req)
	if err != nil {
		z.L().Warn("请求失败", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close() // 必须关闭响应体 [1,5,8](@ref)
	z.L().Debug("请求Github", zap.Any("resp", resp))
	if resp.StatusCode != 200 {
		z.L().Warn("请求失败", zap.Any("resp", resp))
		return nil, fmt.Errorf("请求失败 %v %v", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		z.L().Warn("github请求失败", zap.Error(err))
		return nil, err
	}
	return body, nil
}
func (this *githubApi) Request(githubUser, repoName string) (*model.GitHubModel, error) {
	//this := &GithubApi{}
	body, err := request(githubUser, repoName)
	if err != nil {
		return nil, err
	}
	var result model.GitHubModel
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("github请求失败 %v %v", err, len(body))
	}
	this.proxies = utils.ParseMarkdownCodeToStringArray(result.Body)
	return &result, nil
}

func (this *githubApi) defaultRequest() (*model.GitHubModel, error) {
	if this.userName == "" {
		return nil, errors.New("请指定github的用户名")
	}
	if this.repoName == "" {
		return nil, errors.New("请指定github的仓库名")
	}

	return this.Request(this.userName, this.repoName)
}

func (this *githubApi) CheckUpgrade(fullName string) (map[string]interface{}, error) {
	if fullName == "" {
		return nil, errors.New("fullName is empty")
	}
	r, e := this.defaultRequest()
	if e != nil {
		return nil, e
	}
	oldVersion := utils.GetVersionByFileName(fullName)
	z.L().Debug("升级", zap.String("最新版本", r.TagName), zap.String("本地版本", oldVersion))
	hasNewVersion := utils.CompareVersions(r.TagName, oldVersion)
	z.L().Debug("计算结果", zap.Int("hasNewVersion", hasNewVersion))

	var fullAsset, patchAsset model.Assets
	if hasNewVersion > 0 {
		newVersionAppName := utils.ReplaceNewVersionBinName(fullName, r.TagName)
		var fullUrl, patchUrl string
		var fullSize, patchSize string
		patchName := fmt.Sprintf("%s.patch", newVersionAppName)
		for _, asset := range r.Assets {
			if strings.Compare(strings.ToLower(asset.Name), strings.ToLower(newVersionAppName)) == 0 {
				fullUrl = asset.BrowserDownloadUrl
				fullSize = utils.ByteCountIEC(uint64(asset.Size))
				fullAsset = asset
			} else if strings.Compare(strings.ToLower(asset.Name), strings.ToLower(patchName)) == 0 {
				patchUrl = asset.BrowserDownloadUrl
				patchSize = utils.ByteCountIEC(uint64(asset.Size))
				patchAsset = asset
			}
		}

		if hasNewVersion != 1 {
			//版本之间只有相差一个版本号才能差量升级
			patchUrl = ""
		}
		index := strings.Index(r.Body, "---")
		releaseNote := r.Body
		if index > 0 {
			releaseNote = releaseNote[:index]
		}

		text := strings.Builder{}
		if fullSize != "" {
			text.WriteString("\r\n* 全量包：")
			text.WriteString(fullSize)
		}

		if patchSize != "" {
			text.WriteString("\r\n* 差量包：")
			text.WriteString(patchSize)
		}

		return map[string]interface{}{
			"raw":          *r,
			"fullAsset":    fullAsset,
			"patchAsset":   patchAsset,
			"fullSize":     fullAsset.Size,
			"patchSize":    patchAsset.Size,
			"fullUrl":      fullUrl,
			"patchUrl":     patchUrl,
			"publishedAt":  r.PublishedAt,
			"htmlUrl":      r.HtmlUrl,
			"version":      r.TagName,
			"releaseNotes": fmt.Sprintf("### ✅ 新版本\r\n* %s%s\r\n%s", r.TagName, text.String(), releaseNote),
		}, nil
	} else {
		return nil, fmt.Errorf("已是最新版本～")
	}
}

func (this *githubApi) GetProxyUrls(fileUrl string) []string {
	newProxy := make([]string, 0)
	if this.proxies == nil || len(this.proxies) <= 0 {
		newProxy = append(newProxy, fileUrl)
	} else {
		for _, proxy := range this.proxies {
			newUrl := fmt.Sprintf("%s%s", proxy, fileUrl)
			newProxy = append(newProxy, newUrl)
		}
		newProxy = append(newProxy, fileUrl)
	}
	return newProxy
}

func (this *githubApi) GetModel() *model.GitHubModel {
	r, e := this.defaultRequest()
	if e != nil {
		z.L().Warn("GetModel", zap.Error(e))
		return nil
	}
	return r
}

func (this *githubApi) SetName(userName, repoName string) *githubApi {
	if userName != "" {
		this.userName = userName
	}
	if repoName != "" {
		this.repoName = repoName
	}
	return this
}
func (this *githubApi) GetDownloadUrl(fn func(string, *model.Assets) bool) string {
	r, e := this.defaultRequest()
	if e != nil {
		z.L().Warn("GetDownloadUrl", zap.Error(e))
		return ""
	}
	if r == nil {
		z.L().Warn("this.result is nil")
	} else if r.Assets != nil {
		for _, asset := range r.Assets {
			if fn != nil && fn(r.TagName, &asset) {
				return asset.BrowserDownloadUrl
			}
		}
	}
	return ""
}
func (this *githubApi) GetDownloadUrls(fn func(string, *model.Assets) bool) []string {
	r, e := this.defaultRequest()
	if e != nil {
		z.L().Warn("GetDownloadUrls", zap.Error(e))
		return nil
	}
	if r == nil {
		z.L().Warn("this.result is nil")
		return nil
	} else if r.Assets != nil {
		urls := make([]string, 0)
		for _, asset := range r.Assets {
			if fn != nil && fn(r.TagName, &asset) {
				urls = append(urls, asset.BrowserDownloadUrl)
			}
		}
		return urls
	}
	return nil
}
