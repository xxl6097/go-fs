package gs

import (
	"fmt"

	"github.com/xxl6097/go-service/internal"
	"github.com/xxl6097/go-service/pkg/gs/igs"
)

// InitLog everyType 0：每天，1：每小时，2：每10分钟，3：每分钟 （切割文件）
func InitLog() {
	internal.InitLog()
}

func Run(srv igs.IService) error {
	if srv == nil {
		return fmt.Errorf("请继承igs.IService接口")
	}
	return internal.NewCore(srv).Run()
}
