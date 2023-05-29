package network

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"raft/raft"
	"time"
)

func checkFault(c *gin.Context) {
	ip := c.ClientIP()
	raft.DPrintf("client ip: [%s]", ip)
	var shouldBlock bool
	var delayTime time.Duration
	mu.Lock()
	node, ok := ipToNode[ip]
	if !ok {
		shouldBlock = true
	} else {
		raft.DPrintf("get message from node: %d", node)
		delayTime = nodeDelay[node]
		shouldBlock = nodeBlock[node] || crashed
	}
	mu.Unlock()
	if shouldBlock {
		c.Abort()
		c.String(http.StatusOK, "Packet loss")
		raft.DPrintf("block request from [%d]", node)
		return
	} else {
		time.Sleep(delayTime)
		c.Next()
	}
}
