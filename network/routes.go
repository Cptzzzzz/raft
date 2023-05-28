package network

import "raft/global"

func InitRoutes() {
	client := global.Router.Group("/client")

	client.POST("/operate", operate)
	client.POST("/state", state)
	client.POST("/block", block)
	client.POST("/crash", crash)
	client.POST("/recover", recovery)
	client.POST("/delay", delay)

	node := global.Router.Group("/raft")
	node.Use(checkTerm)
	node.POST("/append-entries", appendEntries)
	node.POST("/request-vote", requestVote)
}
