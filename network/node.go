package network

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"raft/global"
	"raft/raft"
)

func appendEntries(c *gin.Context) {
	var args global.AppendEntriesArgs
	var reply global.AppendEntriesReply
	if err := c.ShouldBindJSON(&args); err != nil {
		c.String(http.StatusBadGateway, "")
	}
	raft.Rf.AppendEntries(&args, &reply)
	c.JSON(http.StatusOK, reply)
}

func requestVote(c *gin.Context) {
	var args global.RequestVoteArgs
	var reply global.RequestVoteReply
	if err := c.ShouldBindJSON(&args); err != nil {
		c.String(http.StatusBadGateway, "")
	}
	raft.Rf.RequestVote(&args, &reply)
	c.JSON(http.StatusOK, reply)
}
