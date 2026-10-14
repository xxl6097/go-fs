package gs

import (
	"encoding/json"
	"net/http"

	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"
)

type GeneralResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
	Raw  []byte `json:"-"`
}

func (g *GeneralResponse) response(code int, msg string, data any) *GeneralResponse {
	g.Msg = msg
	g.Code = code
	g.Data = data
	return g
}
func (g *GeneralResponse) Response(code int, msg string) *GeneralResponse {
	return g.response(code, msg, nil)
}
func (g *GeneralResponse) Result(code int, msg string, data any) *GeneralResponse {
	return g.response(code, msg, data)
}
func (g *GeneralResponse) StatusCode(code int) *GeneralResponse {
	return g.response(code, "", nil)
}
func (g *GeneralResponse) Err(err error) *GeneralResponse {
	//glog.Error(err)
	return g.Response(-1, err.Error())
}
func (g *GeneralResponse) Error(msg string) *GeneralResponse {
	//glog.Error(msg)
	return g.Response(-1, msg)
}
func (g *GeneralResponse) Any(data any) *GeneralResponse {
	return g.response(0, "ok", data)
}
func (g *GeneralResponse) Object(msg string, data any) *GeneralResponse {
	return g.response(0, msg, data)
}
func (g *GeneralResponse) Ok(msg string) *GeneralResponse {
	return g.Response(0, msg)
}
func (g *GeneralResponse) Sucess(msg string, data any) *GeneralResponse {
	return g.response(0, msg, data)
}

func Response(r *http.Request) (*GeneralResponse, func(w http.ResponseWriter)) {
	res := &GeneralResponse{Code: 0}
	return res, func(w http.ResponseWriter) {
		defer func() {
			if res.Code != 0 {
				z.L().Warn("Http response error", zap.String("URL", r.URL.Path), zap.Any("res", res))
			}
		}()

		w.WriteHeader(200)
		var data []byte
		if res.Data == nil {
			//res.Data = utils.GetTime()
		}
		if res.Raw != nil {
			data = res.Raw
			z.L().Info("Http response", zap.String("Method", r.Method), zap.String("URL", r.URL.Path), zap.String("raw", string(res.Raw)))
		} else {
			//glog.Infof("Http response [%s %s]: res: %v", r.Method, r.URL.Path, res)
			bb, err := json.Marshal(res)
			if err != nil {
				z.L().Error("marshal error", zap.Error(err))
				w.WriteHeader(400)
				return
			}
			data = bb
		}
		if len(data) > 0 {
			_, _ = w.Write(data)
		}
	}
}
