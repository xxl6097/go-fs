package utils

import (
	"fmt"
	"sync"
	"time"
)

var instance *IDGenerator
var once sync.Once

func GetInstance() *IDGenerator {
	once.Do(func() {
		instance = NewIDGenerator(1000)
	})
	return instance
}

func GetInt64ID() int64 {
	return GetInstance().NextID()
}

func GetID() string {
	return fmt.Sprintf("%d", GetInstance().NextID())
}

// IDGenerator 这个 ID 生成器实现了类似 Snowflake 算法的机制，生成 64 位唯一 ID，结构为：
//
// 41 位时间戳（精确到毫秒，支持约 69 年）
// 10 位机器 ID（最多支持 1024 台机器）
// 12 位序列号（每毫秒内最多生成 4096 个 ID）
//
// 特点：
//
// 生成的 ID 趋势递增
// 分布式环境下保证唯一性（需分配不同机器 ID）
// 单机每秒可生成约 400 万个 ID
// 时间回拨问题处理：当检测到时间回拨时，会等待时间恢复正常
//
// 使用时，只需初始化生成器并调用 NextID () 方法即可获取唯一 ID。注意要为不同机器分配不同的 machineID（0-1023 之间的值）。
//
//	结构体用于生成唯一ID
type IDGenerator struct {
	mu        sync.Mutex
	lastTime  int64
	sequence  int32
	machineID int32
}

// NewIDGenerator 创建一个新的ID生成器实例
func NewIDGenerator(machineID int32) *IDGenerator {
	return &IDGenerator{
		lastTime:  time.Now().UnixNano() / 1e6, // 转换为毫秒
		sequence:  0,
		machineID: machineID,
	}
}

// NextID 生成下一个唯一ID
// ID结构: 41位时间戳 + 10位机器ID + 12位序列号
func (g *IDGenerator) NextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UnixNano() / 1e6 // 毫秒时间戳

	if now == g.lastTime {
		g.sequence = (g.sequence + 1) & 0xFFF // 序列号范围0-4095
		if g.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			for now <= g.lastTime {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		g.sequence = 0
	}

	g.lastTime = now

	// 组装ID: 时间戳部分 | 机器ID部分 | 序列号部分
	id := ((now - 1609459200000) << 22) | // 减去一个固定时间戳，扩展可用年限
		(int64(g.machineID) << 12) |
		int64(g.sequence)

	return id
}
