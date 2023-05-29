package main

import (
	"fmt"
	"math/rand"
	"raft/client/lib"
	"raft/global"
	"time"
)

func main() {
	fmt.Println("----- Delay Testcase -----")
	lib.ReadConfig()
	lib.InitNetwork()
	leader := lib.OneLeader()
	fmt.Println("first leader id", leader)
	for i := 1; i <= 10; i++ {
		delayer := lib.GenerateDelay(-1)
		delayer2 := (delayer + 1 + rand.Int()%4) % 5
		lib.GenerateDelay(delayer2)
		fmt.Println("set delay on node", delayer, delayer2)
		for i := 0; i < 100; i++ {
			lib.SendOperation(leader, global.Command{
				Operator: lib.RandomOperator(),
				Key:      fmt.Sprintf("raft-kv-%d", i%25),
				Value:    fmt.Sprintf("%d", i),
			})
		}
		fmt.Println("waiting to reach consensus...")
		time.Sleep(time.Second * 5)
		reply := lib.GetResult()
		if reply.Message == "ok" {
			fmt.Println("passed loop", i)
			fmt.Println("log length", len(reply.Data[0]))
		} else {
			panic("unable to reach consensus")
		}
		lib.RecoverDelay(delayer)
		lib.RecoverDelay(delayer2)
		fmt.Println("remove delay on node", delayer, delayer2)
		oldLeader := leader
		fmt.Println("leader", leader, "crash down")
		lib.GenerateCrash(oldLeader)
		fmt.Println("waiting for a new leader...")
		time.Sleep(2 * time.Second)
		lib.RecoverCrash(oldLeader)
		leader = lib.OneLeader()
		fmt.Println("new leader is", leader)
	}
	lib.ClearJudge()
}
