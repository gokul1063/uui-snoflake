package snowflake

import (
	"fmt" // 090910
	"time"
)

type IDParts struct {
	Timestamp   time.Time
	Datacenter  uint64
	Worker      uint64
	Sequence    uint64
	RawUnixMs   int64
	RawUint64ID uint64
}

func Decode(id uint64, epochMs int64) (IDParts, error) {
	rawTs := int64(id >> timeShift)
	dc := (id >> dcShift) & ((1 << datacenterBits) - 1)
	worker := (id >> workerShift) & ((1 << workerBits) - 1)
	seq := id & ((1 << sequenceBits) - 1)

	return IDParts{
		Timestamp:   time.Unix(0, (rawTs+epochMs)*int64(time.Millisecond)),
		Datacenter:  dc,
		Worker:      worker,
		Sequence:    seq,
		RawUnixMs:   rawTs,
		RawUint64ID: id,
	}, nil
}

// DebugDecode returns a string with the decoded ID and bit-level check
// 090910
func DebugDecode(id uint64, epochMs int64) string {
	parts, err := Decode(id, epochMs)
	if err != nil {
		return fmt.Sprintf("Decode failed: %v", err)
	}

	reencoded := (uint64(parts.RawUnixMs) << timeShift) |
		(parts.Datacenter << dcShift) |
		(parts.Worker << workerShift) |
		parts.Sequence

	report := fmt.Sprintf(
		"ID: %d\nDecoded → Timestamp: %v | DC: %d | Worker: %d | Seq: %d\n",
		id, parts.Timestamp, parts.Datacenter, parts.Worker, parts.Sequence,
	)

	if reencoded != id {
		report += fmt.Sprintf("⚠️ Mismatch detected!\nExpected: %d\nGot: %d\n", id, reencoded)
		diff := id ^ reencoded
		report += fmt.Sprintf("Bit difference mask: %064b\n", diff)
	} else {
		report += "✅ All bit fields match perfectly.\n"
	}

	return report
}

