package snowflake

import (
	"sync"
	"time"
)

// 推荐阅读：https://www.luozhiyun.com/archives/527
const (
	epoch         = int64(1711900800000) // 起始时间 2024-04-01 00:00:00
	timestampBits = uint(41)             // 41 位时间戳
	nodeBits      = uint(10)             // 10 位节点标识
	sequenceBits  = uint(12)             // 12 位序列号
	timestampMax  = int64(-1) ^ (-1 << timestampBits)
	nodeMax       = int64(-1) ^ (-1 << nodeBits)
	sequenceMax   = int64(-1) ^ (-1 << sequenceBits)

	timestampShift = sequenceBits + nodeBits
	nodeShift      = sequenceBits
)

type SnowFlake struct {
	lock      sync.Mutex
	nodeId    int64
	timestamp int64
	sequence  int64
}

type Option func(*SnowFlake)

func WithEpoch(epoch int64) Option {
	return func(flake *SnowFlake) {
		flake.timestamp = epoch
	}
}

func NewSnowFlake(nodeId int64, opts ...Option) *SnowFlake {
	svc := &SnowFlake{
		nodeId:    nodeId,
		timestamp: epoch,
		sequence:  0,
	}
	for _, opt := range opts {
		opt(svc)
	}
	return svc
}

func (s *SnowFlake) Generate() int64 {
	s.lock.Lock()
	defer s.lock.Unlock()
	now := time.Now().UnixMilli()
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMax
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		panic("clock is moving backwards")
	}
	s.timestamp = now
	return (t << timestampShift) | (s.nodeId << nodeShift) | s.sequence
}
