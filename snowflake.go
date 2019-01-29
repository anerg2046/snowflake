package snowflake

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// Epoch 开始时间戳，固定小于当前时间的毫秒，用于缩减时间位数
	Epoch int64 = 1532620800000
	// MaxMachineBit 最大机器标识位数，8->255 10->1023
	MaxMachineBit uint8 = 8
	// MaxSequenceBit 最大序号数位数，12->4095
	MaxSequenceBit uint8 = 12
	// maxMachineID 最大机器标识号
	maxMachineID uint16 = (1 << MaxMachineBit) - 1
	// timestamp 毫秒时间戳
	timestamp int64
	// sequence 每毫秒内的队列
	sequence int64
)

type Node struct {
	mu        sync.Mutex
	machineID uint16
}
type ID int64

func NewNode(machineID uint16) (*Node, error) {
	maxMachineID = (1 << MaxMachineBit) - 1
	if machineID > maxMachineID {
		return nil, errors.New("Machine ID must be between 0 and " + strconv.FormatInt(int64(maxMachineID), 10))
	}
	return &Node{machineID: machineID}, nil
}

func (n *Node) NextID() ID {
	n.mu.Lock()

	now := time.Now().UnixNano() / 1000000
	if now == timestamp {
		sequence++
	} else {
		timestamp = now
		sequence = 0
	}

	timeDiff := now - Epoch
	var nextID ID
	if n.machineID == 0 {
		timeLeftShift := MaxSequenceBit
		nextID = ID((timeDiff << timeLeftShift) | sequence)
	} else {
		timeLeftShift := MaxSequenceBit + MaxMachineBit
		nextID = ID((timeDiff << timeLeftShift) | (int64(n.machineID) << MaxSequenceBit) | sequence)
	}

	n.mu.Unlock()
	return nextID
}

func (n *Node) ParseID(id int64) map[string]interface{} {
	timeLeftShift := MaxSequenceBit
	if n.machineID > 0 {
		timeLeftShift += MaxMachineBit
	}
	timeDiff := id >> timeLeftShift
	idTime := timeDiff + Epoch
	rsp := make(map[string]interface{})
	rsp["timestamp"] = idTime
	rsp["machineID"] = n.machineID
	return rsp
}

func (i ID) Int64() int64 {
	return int64(i)
}

func (i ID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i ID) Base36() string {
	return strings.ToUpper(strconv.FormatInt(int64(i), 36))
}
