package gs

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/glog/pkg/zutil"
	"github.com/xxl6097/go-service/pkg/github"
	"github.com/xxl6097/go-service/pkg/gs/igs"
	"github.com/xxl6097/go-service/pkg/utils"
	"github.com/xxl6097/go-service/pkg/utils/util"
	"go.uber.org/zap"
)

var pool = &sync.Pool{
	New: func() interface{} { return make([]byte, 32*1024) },
}

func update(srv igs.Service, w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	if srv == nil {
		res.Error("srv is nil")
		return
	}
	ctx := r.Context()
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	updir := zutil.AppHome()
	_, _, free, _ := util.GetDiskUsage(updir)
	if free < utils.GetSelfSize()*2 {
		if err := utils.ClearTemp(); err != nil {
			z.L().Warn("/tmp清空失败", zap.Error(err))
		} else {
			z.L().Debug("/tmp清空完成")
		}
	}

	var newFilePath string
	switch r.Method {
	case "PUT", "put":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			res.Response(400, fmt.Sprintf("read request body error: %v", err))
			z.L().Warn("PUT失败", zap.Any("res", res))
			return
		}
		if len(body) == 0 {
			res.Response(400, "升级URL空的哦～")
			z.L().Warn("PUT失败", zap.Any("res", res))
			return
		}
		binUrl := string(body)
		z.L().Debug("升级URL地址", zap.String("binUrl", binUrl))
		newUrl := utils.DownloadFileWithCancelByUrls(github.Api().GetProxyUrls(binUrl), nil)
		newFilePath = newUrl
		break
	case "POST", "post":
		// 获取上传的文件
		file, handler, err := r.FormFile("file")
		if err != nil {
			res.Error("body no file")
			return
		}
		defer file.Close()
		dstFilePath := filepath.Join(zutil.AppHome("temp", "upgrade"), handler.Filename)
		//dstFilePath 名称为上传文件的原始名称
		dst, err := os.Create(dstFilePath)
		if err != nil {
			res.Error(fmt.Sprintf("create file %s error: %v", handler.Filename, err))
			return
		}
		buf := pool.Get().([]byte)
		defer pool.Put(buf)
		_, err = io.CopyBuffer(dst, file, buf)
		_ = dst.Close()
		if err != nil {
			res.Error(err.Error())
			return
		}
		newFilePath = dstFilePath
		break
	default:
		res.Error("位置请求方法")
	}
	if newFilePath != "" {
		z.L().Sugar().Info("开始升级 ", newFilePath)
		err := srv.Upgrade(ctx, newFilePath)
		z.L().Warn("升级结果", zap.Error(err))
		if err == nil {
			res.Ok("升级成功～")
		} else {
			res.Error(fmt.Sprintf("更新失败～%v", err))
		}

	}
}

func checkVersion(name string, w http.ResponseWriter, r *http.Request) {
	res, f := Response(r)
	defer f(w)
	if name == "" {
		name = zutil.AppName()
	}
	data, err := github.Api().CheckUpgrade(name)
	if err != nil {
		res.Err(err)
	} else {
		z.L().Debug("version", zap.Any("data", data))
		res.Any(data)
	}
}

func ApiCheckVersion(binName string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		checkVersion(binName, w, r)
	}
}

func ApiUpdate(srv igs.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		update(srv, w, r)
	}
}
