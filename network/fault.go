package network

import (
	"raft/global"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var nodeDelay []time.Duration
var nodeBlock []bool
var ipToNode map[string]int
var crashed bool

func InitFault() {
	nodeDelay = make([]time.Duration, len(global.Peers))
	nodeBlock = make([]bool, len(global.Peers))
	ipToNode = make(map[string]int)
	for i := len(global.Peers) - 1; i >= 0; i-- {
		nodeDelay[i] = time.Duration(0)
		nodeBlock[i] = false
	}
	for index, val := range global.Peers {
		ipToNode[strings.Split(val, ":")[0]] = index
	}
	crashed = false
}
