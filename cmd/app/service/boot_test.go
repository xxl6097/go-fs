package service

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/xxl6097/gofs"
)

// legacyConfig 只保留 v0.0.110 配置里与账号锁定相关的字段。
// gob 按字段名匹配，旧版本流里没有的新字段在解码后应保持默认值。
type legacyConfig struct {
	AuthFailLimit  int
	AuthFailWindow time.Duration
}

func TestDecodeConfigKeepsNewFieldDefaults(t *testing.T) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(legacyConfig{
		AuthFailLimit:  10,              // 旧默认值
		AuthFailWindow: 5 * time.Minute, // 旧默认值
	})
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	cfg, err := decodeConfig(buf.Bytes())
	if err != nil {
		t.Fatalf("decodeConfig: %v", err)
	}

	// 流里存在的字段仍按旧值覆盖默认值。
	if cfg.AuthFailLimit != 10 {
		t.Errorf("AuthFailLimit = %d，期望 10", cfg.AuthFailLimit)
	}
	if cfg.AuthFailWindow != 5*time.Minute {
		t.Errorf("AuthFailWindow = %v，期望 5m", cfg.AuthFailWindow)
	}
	// 流里缺失的新字段保留默认值，而不是零值。
	def := gofs.Default()
	if cfg.AuthLockMaxDelay != def.AuthLockMaxDelay {
		t.Errorf("AuthLockMaxDelay = %v，期望默认 %v", cfg.AuthLockMaxDelay, def.AuthLockMaxDelay)
	}
	if len(cfg.AuthLockExempt) != len(def.AuthLockExempt) {
		t.Errorf("AuthLockExempt = %v，期望默认 %v", cfg.AuthLockExempt, def.AuthLockExempt)
	}
	// 启动路径真正会走到的校验。
	if err := cfg.Normalize(); err != nil {
		t.Errorf("Normalize: %v", err)
	}
}
