package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"unicode"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"go.uber.org/zap"
)

func CheckFileOrDownload(ctx context.Context, fileUrlOrLocal string) (string, error) {
	if IsURL(fileUrlOrLocal) {
		filePath, err := DownloadWithCancel(ctx, fileUrlOrLocal, nil)
		if err != nil {
			z.L().Error("下载失败", zap.String("fileUrlOrLocal", fileUrlOrLocal), zap.Error(err))
			return "", err
		}
		z.L().Sugar().Debug("下载成功", filePath)
		return filePath, nil
	} else if FileExists(fileUrlOrLocal) {
		z.L().Sugar().Debug("检测为本地文件", fileUrlOrLocal)
		return fileUrlOrLocal, nil
	} else {
		z.L().Sugar().Error("无法识别的文件", fileUrlOrLocal)
		return "", errors.New("无法识别的文件" + fileUrlOrLocal)
	}
}

func FileSize(file string) {
	// 获取文件信息
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	// 获取文件大小（字节）
	size := fileInfo.Size()

	z.L().Debug(fmt.Sprintf("文件 %s", file))
	z.L().Debug(fmt.Sprintf("文件大小: %d 字节", size))
	z.L().Debug(fmt.Sprintf("格式化显示: %.2f KB", float64(size)/1024))
	z.L().Debug(fmt.Sprintf("格式化显示: %.2f MB", float64(size)/(1024*1024)))
}
func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// FileExists 用于判断文件是否存在
func FileExists(filePath string) bool {
	// 调用 os.Stat 函数获取文件信息
	f, err := os.Stat(filePath)
	// 判断是否为文件不存在的错误
	if os.IsNotExist(err) {
		return false
	}
	if f != nil {
		fmt.Println(ByteCountIEC(uint64(f.Size())), filePath)
	}
	// 若有其他错误或无错误，认为文件存在
	return true
}
func CheckDirector(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// 存在，删除
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func ResetDirector(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		// 存在，删除
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
		return os.MkdirAll(path, 0755)
	}
	return os.MkdirAll(path, 0755)
}

func DeleteAllDirector(filePath string) error {
	err := os.RemoveAll(filePath)
	if err != nil {
		msg := fmt.Errorf("删除失败[%v]: %s,%v\n", os.Getpid(), filePath, err)
		z.L().Sugar().Error(msg)
		return msg
	}
	//z.L().Sugar().Infof("删除成功[%v]: %s\n", os.Getpid(), filePath)
	return err
}

//	func CopyToTemp(srcFile string) error {
//		src, err := os.Open(srcFile) // can not use args[0], on Windows call openp2p is ok(=openp2p.exe)
//		if err != nil {
//			fmt.Printf("打开源文件失败：%v\n", err)
//			return err
//		}
//		var fileSize int64
//		var fileName string
//		finfo, err := src.Stat()
//		if err == nil {
//			fileSize = finfo.Size()
//			fileName = finfo.Name()
//		}
//
//		dstFile := filepath.Join(zutil.TempDir(), filepath.Base(srcFile))
//		defer src.Close()
//		//将本程序复制到目标为止，目标文件名称为配置文件的名称
//		dst, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
//		if err != nil {
//			fmt.Printf("创建目标文件失败：%v\n", err)
//			return err
//		}
//		defer dst.Close()
//		sizeB := float64(fileSize) / 1024 / 1024
//		fmt.Printf("正在拷贝%s[大小：%.2fMB]到%s\n", fileName, sizeB, dstFile)
//		_, err = io.Copy(dst, src)
//		if err != nil {
//			fmt.Printf("拷贝文件失败：%v\n", err)
//			return err
//		}
//		return nil
//	}
func Copy(srcFile, dstFile string) error {
	src, err := os.Open(srcFile) // can not use args[0], on Windows call openp2p is ok(=openp2p.exe)
	if err != nil {
		fmt.Printf("copy 打开源文件失败：%v\n", err)
		return err
	}
	var fileSize int64
	var fileName string
	finfo, err := src.Stat()
	if err == nil {
		fileSize = finfo.Size()
		fileName = finfo.Name()
	}
	defer src.Close()
	if e := EnsureDir(dstFile); e != nil {
		fmt.Printf(" EnsureDir：%v %v %v\n", srcFile, dstFile, e)
		return err
	}
	//将本程序复制到目标为止，目标文件名称为配置文件的名称
	dst, err := os.OpenFile(dstFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0775)
	if err != nil {
		fmt.Printf("copy 创建目标文件失败：%v %v %v\n", srcFile, dstFile, err)
		return err
	}
	defer dst.Close()
	sizeB := float64(fileSize) / 1024 / 1024
	fmt.Printf("正在拷贝%s[大小：%.2fMB]到%s\n", fileName, sizeB, dstFile)
	_, err = io.Copy(dst, src)
	if err != nil {
		fmt.Printf("拷贝文件失败：%v\n", err)
		return err
	}
	return nil
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err) // 如果路径不存在或权限不足返回 false
	}
	return true
}
func CleanExt(name string) string {
	filename := path.Base(name) // 获取文件名"app.log"
	nameOnly := filename[:len(filename)-len(path.Ext(filename))]
	return nameOnly
}

// ToUpperFirst 将字符串的首字母转换为大写
func ToUpperFirst(s string) string {
	if s == "" {
		return s
	}
	// 将字符串转换为符文切片
	r := []rune(s)
	// 将首字符转换为大写
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func ClearTemp() error {
	tempDir := zutil.TempDir()
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("读取目录失败: %v", err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(tempDir, entry.Name())
		err = os.RemoveAll(fullPath)
		if err != nil {
			z.L().Sugar().Debugf("删除失败 %s  %v", fullPath, err)
		} else {
			z.L().Sugar().Debugf("删除成功 %s", fullPath)
		}
	}
	return err
}

func EnsureDir(filePath string) error {
	// 获取文件所在的**目录路径**
	dir := filepath.Dir(filePath)

	// MkdirAll 会自动：
	// 1. 目录不存在 → 创建
	// 2. 目录已存在 → 不报错、不覆盖
	// 0755 是常用的目录权限
	return os.MkdirAll(dir, 0755)
}

func EnsureDir1(filePath string) error {
	// 获取文件所在的**目录路径**
	dir := filepath.Dir(filePath)

	// MkdirAll 会自动：
	// 1. 目录不存在 → 创建
	// 2. 目录已存在 → 不报错、不覆盖
	// 0755 是常用的目录权限
	return os.MkdirAll(dir, 0755)
}
