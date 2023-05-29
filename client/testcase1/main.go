package main

import (
	"fmt"
	"raft/client/lib"
	"raft/global"
	"time"
)

func main() {
	fmt.Println("----- Normal Testcase -----")
	lib.ReadConfig()
	lib.InitNetwork()
	leader := lib.OneLeader()
	fmt.Println("leader is ", leader)
	start := time.Now()
	for i := 0; i < 1000; i++ {
		lib.SendOperation(leader, global.Command{
			Operator: lib.RandomOperator(),
			Key:      fmt.Sprintf("raft-kv-%d", i%25),
			Value:    fmt.Sprintf("%d", i),
		})
	}
	time.Sleep(time.Second * 3)
	reply := lib.GetResult()
	if reply.Message == "ok" {
		fmt.Println("success")
		fmt.Println("time used:", time.Now().Sub(start))
		fmt.Println("log length", len(reply.Data[0]))
	} else {
		fmt.Println(reply.Message)
	}
	lib.ClearJudge()
}
