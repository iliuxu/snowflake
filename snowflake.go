package snowflake

import (
	"encoding/base64"
	"errors"
	"strconv"
	"sync"
	"time"
)

var(
	nodeBits uint8 = 10
	stepBits uint8 = 12
	nodeMax int64 = -1 ^(-1 << nodeBits)
	nodeMask int64 = nodeMax << stepBits
	stepMask int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
	epoch int64 = 1540477876017
)

type Node struct {
	mu sync.Mutex
	time int64
	node int64
	step int64
}

type ID int64

func NewNode(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
			return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(nodeMax, 10))
	}
	return &Node{
		time:0,
		node:node,
		step:0,
	},nil
}

func (n *Node) Generate() ID {
	n.mu.Lock()
	now := time.Now().UnixNano() / 1000000
	if n.time == now {
		n.step =(n.step+1) & stepMask
		if n.step == 0 {
			for now<=n.time {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}
	n.time = now
	r := ID((now-epoch) << timeShift | (n.node << nodeShift) | (n.step),)
	n.mu.Unlock()
	return r
}

func (f ID) Int64() int64 {
	return int64(f)
}

func (f ID) Time() int64 {
	return (int64(f) >> timeShift) + epoch
}

func (f ID) Node() int64 {
	return int64(f) & nodeMask >> nodeShift
}

func (f ID) Step() int64 {
	return int64(f) & stepMask
}

func (f ID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

func (f ID) Bytes() []byte {
	return []byte(f.String())
}

func (f ID) Base64() string {
	return base64.StdEncoding.EncodeToString(f.Bytes())
}

func (f ID) Base2() string {
	return strconv.FormatInt(int64(f), 2)
}
