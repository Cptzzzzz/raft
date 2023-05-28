package network

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"raft/global"
	"raft/raft"
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

}

func block(c *gin.Context) {

}

func crash(c *gin.Context) {

}

func recovery(c *gin.Context) {

}

func delay(c *gin.Context) {

}
