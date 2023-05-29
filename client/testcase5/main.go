package main

import (
	"fmt"
	"math/rand"
	"raft/client/lib"
	"raft/global"
	"time"
)

func main() {
	fmt.Println("----- Multi Testcase -----")
	lib.ReadConfig()
	lib.InitNetwork()
	leader := lib.OneLeader()
	fmt.Println("The first leader is", leader)
	for i := 1; i <= 10; i++ {
		write(10)
		randomEvent()
		write(10)
		checkConsensus(i)
	}
}

func randomEvent() {
	event := rand.Int() % 3
	switch event {
	case 0:
		crashAndDelay()
	case 1:
		partitionAndDelay()
	default:
		changeLeader()
	}
}

func changeLeader() {
	oldLeader := lib.OneLeader()
	fmt.Printf("--- Change Leader %d ---\n", oldLeader)
	lib.GenerateCrash(oldLeader)
	time.Sleep(3 * time.Second)
	write(30)
	lib.RecoverCrash(oldLeader)

}
func crashAndDelay() {
	fmt.Println("--- Crash and Delay ---")
	who := rand.Int() % 8
	leader := lib.OneLeader()
	if who >= 4 {
		who = leader
	} else {
		who = (leader + who) % 5
	}
	delayer1 := lib.GenerateDelay(-1)
	delayer2 := (delayer1 + 1 + rand.Int()%4) % 5
	fmt.Printf("leader: [%d]. delayer: [%d],[%d]. crasher: [%d]\n", leader, delayer1, delayer2, who)
	lib.GenerateDelay(delayer2)
	lib.GenerateCrash(who)
	fmt.Println("Waiting for the situation......")
	time.Sleep(3 * time.Second)
	write(30)
	lib.RecoverDelay(delayer1)
	lib.RecoverDelay(delayer2)
	lib.RecoverCrash(who)
}
func partitionAndDelay() {
	fmt.Println("--- Partition and Delay ---")
	who := rand.Int() % 8
	leader := lib.OneLeader()
	if who >= 4 {
		who = leader
	} else {
		who = (leader + who) % 5
	}
	lib.GeneratePartition(who)
	delayer1 := lib.GenerateDelay(-1)
	delayer2 := (delayer1 + 1 + rand.Int()%4) % 5
	fmt.Printf("leader: [%d]. delayer: [%d],[%d].\n", leader, delayer1, delayer2)
	fmt.Println("Waiting for the situation......")
	time.Sleep(3 * time.Second)
	write(30)
	lib.RecoverPartition()
	lib.RecoverDelay(delayer1)
	lib.RecoverDelay(delayer2)
}

func write(number int) {
	leader := lib.OneLeader()
	fmt.Printf("Write %d Logs to Leader %d\n", number, leader)
	for i := 0; i < number; i++ {
		lib.SendOperation(leader, global.Command{
			Operator: lib.RandomOperator(),
			Key:      fmt.Sprintf("raft-kv-%d", rand.Int()%100),
			Value:    fmt.Sprintf("%d", rand.Int()),
		})
	}
}

func checkConsensus(pass int) {
	fmt.Println("Waiting for consensus......")
	time.Sleep(3 * time.Second)
	reply := lib.GetResult()
	if reply.Message == "ok" {
		fmt.Println("Passed loop:", pass, "Log length:", len(reply.Data[0]))
	} else {
		panic("unable to reach consensus")
	}
}
