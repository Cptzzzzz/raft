package main

import (
	"fmt"
	"raft/client/lib"
	"raft/global"
	"time"
)

func main() {
	fmt.Println("----- Partition Testcase -----")
	lib.ReadConfig()
	lib.InitNetwork()
	leader := lib.OneLeader()
	fmt.Println("first leader id", leader)
	for i := 1; i <= 10; i++ {
		lib.GeneratePartition(leader)
		fmt.Println("waiting for new leader...")
		time.Sleep(2 * time.Second)
		leader = lib.OneLeader()
		fmt.Println("new leader is", leader)
		lib.RecoverPartition()
		for j := 0; j < 20; j++ {
			lib.SendOperation(leader, global.Command{
				Operator: lib.RandomOperator(),
				Key:      fmt.Sprintf("raft-kv-%d", j%25),
				Value:    fmt.Sprintf("%d", j),
			})
		}
		fmt.Println("waiting to reach consensus...")
		time.Sleep(3 * time.Second)
		reply := lib.GetResult()
		if reply.Message == "ok" {
			fmt.Println("passed loop", i)
			fmt.Println("log length", len(reply.Data[0]))
		} else {
			panic("unable to reach consensus")
		}
	}
	lib.ClearJudge()
}
