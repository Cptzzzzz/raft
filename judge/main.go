package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"raft/global"
	"reflect"
	"sync"
)

var data []map[int]global.Command
var mu sync.Mutex

type Reply struct {
	Data    map[int][]global.Command `json:"data"`
	Message string                   `json:"message"`
}

func main() {
	g := gin.Default()
	g.POST("/msg", solve)
	g.POST("/reset", reset)
	g.POST("/result", result)
	data = make([]map[int]global.Command, 5)
	for i := 0; i < 5; i++ {
		data[i] = make(map[int]global.Command)
	}
	g.Run("0.0.0.0:80")
}

func solve(c *gin.Context) {
	var args global.JudgeArgs
	if err := c.ShouldBindJSON(&args); err != nil {
		c.String(http.StatusOK, "")
	}
	mu.Lock()
	defer mu.Unlock()
	data[args.Peer][args.CommandIndex] = args.Command
}

func reset(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	for i := 0; i < 5; i++ {
		data[i] = make(map[int]global.Command)
	}
	c.String(http.StatusOK, "ok")
}

func result(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	reply := Reply{}
	res := make(map[int][]global.Command)
	for i := 0; i < 5; i++ {
		res[i] = make([]global.Command, len(data[i]))
		for k, v := range data[i] {
			if k >= len(data[i]) {
				reply.Message = "log length wrong"
				c.JSON(http.StatusBadGateway, reply)
				return
			}
			res[i][k] = v
		}
	}
	reply.Data = res
	for i := 0; i < 4; i++ {
		if !reflect.DeepEqual(res[i], res[4]) {
			reply.Message = "log content wrong"
			c.JSON(http.StatusBadGateway, reply)
			return
		}
	}
	reply.Message = "ok"
	c.JSON(http.StatusOK, reply)
}
