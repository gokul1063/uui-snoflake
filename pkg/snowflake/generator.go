package snowflake

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Bits allocation
const (
	timestampBits  uint8 = 41
	datacenterBits uint8 = 5
	workerBits     uint8 = 5
	sequenceBits   uint8 = 12

	// Shifts allocation
	sequenceMask = (1 << sequenceBits) - 1
	workerShift  = sequenceBits
	dcShift      = sequenceBits + workerBits
	timeShift    = sequenceBits + workerBits + datacenterBits
)

// Constants
const defaultEpoch int64 = 1735689600000
var ErrInvalidNode = errors.New("Datacenter/worker id out of range.")

// Structures
type GeneratorConfig struct {
	Epoch        int64
	DatacenterId uint64
	WorkerId     uint64
	Shards       int
}

type shard struct {
	mu           sync.Mutex
	lastUnixMs   uint64
	sequence     uint64
	DatacenterId uint64
	WorkerId     uint64
}

type Generator struct {
	epochMs int64
	Shards  []*shard
	nShards uint32
	rr      uint32
}

// Generator creation
func NewGenerator(cfg GeneratorConfig) (*Generator, error) {
	if cfg.Epoch == 0 {
		cfg.Epoch = defaultEpoch
	}

	maxDc := uint64((1 << datacenterBits) - 1)
	maxWorker := uint64((1 << workerBits) - 1)

	if cfg.DatacenterId > maxDc || cfg.WorkerId > maxWorker {
		return nil, ErrInvalidNode
	}

	if cfg.Shards <= 0 {
		cfg.Shards = runtime.NumCPU() * 2
	}

	g := &Generator{
		epochMs: cfg.Epoch,
		Shards:  make([]*shard, cfg.Shards),
		nShards: uint32(cfg.Shards),
	}

	for i := 0; i < cfg.Shards; i++ {
		g.Shards[i] = &shard{
			DatacenterId: cfg.DatacenterId,
			WorkerId:     cfg.WorkerId,
			lastUnixMs:   0,
			sequence:     0,
		}
	}

	return g, nil
}

// Current time in ms since epoch
func (g *Generator) nowMs() uint64 {
	return uint64(time.Now().UnixNano()/1e6 - g.epochMs)
}

// Generate next unique ID
func (g *Generator) Next() (uint64, error) {
	idx := int(atomic.AddUint32(&g.rr, 1) % g.nShards)
	s := g.Shards[idx]

	s.mu.Lock()
	defer s.mu.Unlock()

	now := g.nowMs()

	if now < s.lastUnixMs {
		diff := s.lastUnixMs - now
		if diff > 5 {
			return 0, errors.New("clock moved backwards too far")
		}
		for now < s.lastUnixMs {
			now = g.nowMs()
		}
	}

	if now == s.lastUnixMs {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.lastUnixMs {
				now = g.nowMs()
			}
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}

	s.lastUnixMs = now

	id := (uint64(now) << uint64(timeShift)) |
		(uint64(s.DatacenterId) << uint64(dcShift)) |
		(uint64(s.WorkerId) << uint64(workerShift)) |
		(s.sequence)

	return id, nil
}

