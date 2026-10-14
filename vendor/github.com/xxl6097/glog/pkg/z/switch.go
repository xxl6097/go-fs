package z

import (
	"sync/atomic"

	"go.uber.org/zap/zapcore"
)

// logEnabled 日志总开关，默认开启。
// 可通过 SwitchLog / EnableLog / DisableLog 在运行时动态控制，
// 关闭后所有级别的日志（控制台、主文件、错误文件）都会被丢弃。
var logEnabled atomic.Bool

func init() {
	logEnabled.Store(true)
}

// SwitchLog 设置日志总开关。
//
//	SwitchLog(true)  开启日志
//	SwitchLog(false) 关闭日志
func SwitchLog(enable bool) {
	logEnabled.Store(enable)
}

// EnableLog 开启日志。
func EnableLog() {
	logEnabled.Store(true)
}

// DisableLog 关闭日志。
func DisableLog() {
	logEnabled.Store(false)
}

// IsLogEnabled 返回日志开关当前状态，true 表示开启。
func IsLogEnabled() bool {
	return logEnabled.Load()
}

// switchCore 包装 zapcore.Core，在开关关闭时丢弃所有日志。
// 它包裹整个 tee 核心，因此能同时作用于控制台、主文件与错误文件输出，
// 也能作用于通过 zap.L() / L() 获取的全局 Logger 以及包级 oldLog。
type switchCore struct {
	zapcore.Core
	enabled *atomic.Bool
}

var _ zapcore.Core = (*switchCore)(nil)

// Enabled 关闭时对所有级别返回 false，避免下游把日志写入 CheckedEntry。
func (s *switchCore) Enabled(lvl zapcore.Level) bool {
	return s.enabled.Load() && s.Core.Enabled(lvl)
}

// With 派生新的核心时，继续沿用同一个开关状态。
func (s *switchCore) With(fields []zapcore.Field) zapcore.Core {
	return &switchCore{Core: s.Core.With(fields), enabled: s.enabled}
}

// Check 关闭时原样返回 ce（即 nil），从而跳过写入与 Hook。
func (s *switchCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if !s.enabled.Load() {
		return ce
	}
	return s.Core.Check(ent, ce)
}
