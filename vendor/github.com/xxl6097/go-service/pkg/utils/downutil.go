package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"go.uber.org/zap"
)

type Progress struct {
	Loaded int64 `json:"loaded"`
	Total  int64 `json:"total"`
}

// DownloadFileWithCancelByUrls 并发地从多个 url 下载同一个文件，谁先下载成功就用谁，
// 并立即取消其余仍在进行的下载。返回成功下载到本地的文件路径；全部失败时返回空字符串。
//
// 这里刻意不走 DynamicSelect：DynamicSelect 会把每个 goroutine 的返回值都投递到结果
// channel，下载失败的 goroutine 只能返回空路径，一旦这个空路径先被读到，上层就会把
// 「下载成功但路径为空」当成结果，导致升级被静默跳过。因此成功路径是唯一出口，
// 失败的 goroutine 只打日志、直接返回。
func DownloadFileWithCancelByUrls(urls []string, onProgress func(Progress)) string {
	if len(urls) == 0 {
		return ""
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	success := make(chan string, len(urls)) // 只投递成功的下载结果
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			dstFilePath, err := DownloadWithCancel(ctx, url, onProgress)
			if err != nil {
				// 被其它 url 抢先完成后取消属于正常情况，不必报错
				if !errors.Is(err, context.Canceled) {
					z.L().Warn("下载失败", zap.String("url", url), zap.Error(err))
				}
				return
			}
			FileSize(dstFilePath)
			z.L().Debug("下载成功", zap.String("本地路径", dstFilePath), zap.String("url", url))
			select {
			case success <- dstFilePath:
			case <-ctx.Done():
			}
		}(u)
	}

	// 所有下载都结束时关闭 done，用于「全部失败」时能及时返回，而不是永久阻塞。
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case dst := <-success:
		// 立刻取消其余下载并返回：不等待它们收尾，避免个别地址取消慢（如 DNS 卡住）
		// 拖住整个升级流程；它们会在协程里自行清理各自的临时文件。
		cancel()
		return dst
	case <-done:
		// 全部结束，再确认一次是否有人成功（select 可能随机选中 done，而成功结果已在
		// 缓冲 channel 里）。
		select {
		case dst := <-success:
			return dst
		default:
			z.L().Error("所有下载地址均失败", zap.Int("urls", len(urls)))
			return ""
		}
	}
}

func DownloadWithCancel(ctx context.Context, url string, onProgress func(Progress), args ...string) (string, error) {
	// 创建可取消的 HTTP 请求
	//glog.Debug("开始下载", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// 创建 HTTP 客户端
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	defer resp.Body.Close()
	var dstFile string
	if args != nil && len(args) > 0 {
		dstFile = args[0]
	}
	tempFolder := GetID() //fmt.Sprintf("%d", time.Now().UnixNano())
	if dstFile == "" {
		dstName := GetFileNameFromUrl(url)
		if dstName == "" {
			dstName = GetFilenameFromHeader(resp.Header)
		}
		if dstName == "" {
			fileName := time.Now().Unix()
			dstName = fmt.Sprintf("%d", fileName)
		}
		if dstName != "" {
			dstFile = filepath.Join(zutil.AppHome("temp", "upgrade"), tempFolder, dstName)
		}
	} else {
		dir, f := filepath.Split(dstFile)
		dstFile = filepath.Join(dir, tempFolder, f)
	}
	dir, _ := filepath.Split(dstFile)
	goroutineId := GetGoroutineID()
	_ = ResetDirector(dir)
	// 创建目标文件
	//fmt.Println("os.Create", dstFile)
	outFile, err := os.Create(dstFile)
	if err != nil {
		_ = DeleteAllDirector(dir)
		return "", err
	}
	defer outFile.Close()
	totalSize := resp.ContentLength
	// 分块读取并写入文件
	buf := make([]byte, 4096) // 4KB 缓冲区
	var preProgress float64 = -3.1
	var written int64 // 累计已写入字节数，不依赖文件句柄（避免对已关闭句柄取 Stat 崩溃）
	for {
		select {
		case <-ctx.Done(): // 检查取消信号
			//fmt.Println("下载已取消:", url)
			_ = outFile.Close()
			_ = DeleteAllDirector(dir)
			return "", ctx.Err()
		default:
			n, err1 := resp.Body.Read(buf)
			if err1 != nil && err1 != io.EOF {
				_ = outFile.Close()
				_ = DeleteAllDirector(dir)
				return "", err1
			}
			if n == 0 {
				_ = outFile.Close()
				// 校验下载完整性：仅当服务端给出了有效 Content-Length 时才比对
				// （差量 .patch 等场景 ContentLength 可能为 0/-1，此时跳过校验）。
				// 不一致多为网络截断或磁盘/内存（OpenWRT 的 tmpfs）写满导致的截断文件，
				// 直接返回截断二进制会在后续试跑时崩溃（SIGSEGV）。
				if totalSize > 0 && written != totalSize {
					_ = DeleteAllDirector(dir)
					return "", fmt.Errorf("下载文件不完整：期望 %d 字节，实际 %d 字节（可能磁盘/内存空间不足或网络中断）", totalSize, written)
				}
				z.Println("文件下载完成：", dstFile)
				return dstFile, nil // 正常完成
			}

			w, e := outFile.Write(buf[:n])
			if e != nil {
				_ = outFile.Close()
				_ = DeleteAllDirector(dir)
				return "", fmt.Errorf("写入升级文件失败（可能空间不足）: %w", e)
			}
			written += int64(w)
			if totalSize > 0 {
				progress := float64(written) / float64(totalSize) * 100
				if progress-preProgress > 3 {
					if onProgress != nil {
						onProgress(Progress{Loaded: written, Total: totalSize})
					}
					fmt.Printf("[%d]总大小: %.2fMB 已下载: %.2fMB 进度: %.2f%%\n", goroutineId, float64(totalSize)/1e6, float64(written)/1e6, progress)
					preProgress = progress
				}
			}
		}
	}
}

func GetFileNameFromUrl(rawURL string) string {
	parsedURL, _ := url.Parse(rawURL)
	// 提取路径部分并获取文件名
	fileName := path.Base(parsedURL.Path)
	//fmt.Println("文件名:", fileName) // 输出: document.pdf
	return fileName
}

func GetFilenameFromHeader(header http.Header) string {
	contentDisposition := header.Get("Content-Disposition")
	parts := strings.Split(contentDisposition, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "filename=") {
			fileName := strings.TrimPrefix(part, "filename=")
			fileName = strings.Trim(fileName, `"`) // 去除双引号
			return fileName
		}
	}
	return ""
}

// GetGoroutineID 用于获取当前协程的ID
func GetGoroutineID() uint64 {
	var buf [64]byte
	// 调用runtime.Stack获取当前协程的栈信息
	n := runtime.Stack(buf[:], false)
	// 解析栈信息以提取协程ID
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	var id uint64
	fmt.Sscanf(idField, "%d", &id)
	return id
}

// IsURL 判断给定的字符串是否是一个有效的URL
func IsURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil {
		return false
	}

	return u.Scheme == "http" || u.Scheme == "https"
}
