package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"raft/global"
	"raft/raft"
	"time"
)

func operate(c *gin.Context) {
	args := global.OperateArgs{}
	if err := c.ShouldBindJSON(&args); err != nil {
		c.String(http.StatusBadGateway, "")
	}
	fmt.Println(args)
	index, term, ok, leader := raft.Rf.Append(global.Command{
		Operator: args.Operator,
		Key:      args.Key,
		Value:    args.Value,
	})
	c.JSON(http.StatusOK, global.OperateReply{
		Index:  index,
		Term:   term,
		Ok:     ok,
		Leader: leader,
	})
}

func state(c *gin.Context) {
	var reply global.StateReply
	raft.Rf.GetState(&reply)
	c.JSON(http.StatusOK, reply)
}

func block(c *gin.Context) {
	var args global.BlockArgs
	if err := c.ShouldBindJSON(&args); err != nil {
		c.String(http.StatusBadGateway, "request error")
		return
	}
	mu.Lock()
	for k, v := range args.Block {
		nodeBlock[k] = v
	}
	mu.Unlock()
	c.String(http.StatusOK, "ok")
}

func crash(c *gin.Context) {
	mu.Lock()
	crashed = true
	mu.Unlock()
	raft.Rf.Crash()
	c.String(http.StatusOK, "ok")
}

func recovery(c *gin.Context) {
	mu.Lock()
	crashed = false
	mu.Unlock()
	raft.Rf.Recover()
	c.String(http.StatusOK, "ok")
}

func delay(c *gin.Context) {
	var args global.DelayArgs
	if err := c.ShouldBindJSON(&args); err != nil {
		c.String(http.StatusBadGateway, "request error")
		return
	}
	mu.Lock()
	for k, v := range args.Delay {
		nodeDelay[k] = time.Duration(v) * time.Millisecond
	}
	mu.Unlock()
	c.String(http.StatusOK, "ok")
}
