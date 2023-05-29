package lib

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"raft/global"
	"time"
)

func GetLeader() (int, bool) {
	term := -1
	leader := -1
	ok := false
	for k, _ := range Hosts {
		reply, success := getState(k)
		if success && reply.State == "Leader    " {
			if reply.CurrentTerm > term {
				term = reply.CurrentTerm
				leader = k
				ok = true
			} else if reply.CurrentTerm == term && term != -1 {
				panic("multi leader")
			}
		}
	}
	return leader, ok
}

func getState(peer int) (global.StateReply, bool) {
	var reply global.StateReply
	data := SendClientRequest(Hosts[peer]+"/client/state", make([]byte, 0))
	if data == nil || len(data) == 0 {
		return global.StateReply{}, false
	}
	json.Unmarshal(data, &reply)
	return reply, true
}

func OneLeader() int {
	for {
		leader, ok := GetLeader()
		if ok {
			return leader
		}
		time.Sleep(time.Second)
	}
}

func SendOperation(peer int, command global.Command) global.OperateReply {
	var reply global.OperateReply
	var args global.OperateArgs
	args.Operator = command.Operator
	args.Key = command.Key
	args.Value = command.Value
	body, _ := json.Marshal(&args)
	data := SendClientRequest(Hosts[peer]+"/client/operate", body)
	if data == nil || len(data) == 0 {
		return global.OperateReply{
			Ok: false,
		}
	}
	json.Unmarshal(data, &reply)
	return reply
}

func GetResult() global.JudgeResultReply {
	var reply global.JudgeResultReply
	data := SendClientRequest(Judge+"/result", make([]byte, 0))
	if data == nil || len(data) == 0 {
		return global.JudgeResultReply{}
	}
	json.Unmarshal(data, &reply)
	return reply
}

func ClearJudge() {
	SendClientRequest(Judge+"/reset", make([]byte, 0))
}

func setBlock(peer int, block []bool) {
	var args global.BlockArgs
	args.Block = block
	body, _ := json.Marshal(&args)
	for {
		res := SendClientRequest(Hosts[peer]+"/client/block", body)
		if string(res) == "ok" {
			break
		}
	}
}

func shouldBlock(i, j, peer1, peer2 int) bool {
	if i == peer1 || i == peer2 {
		if j != peer1 && j != peer2 {
			return true
		} else {
			return false
		}
	} else {
		if j != peer1 && j != peer2 {
			return false
		} else {
			return true
		}
	}
}
func GeneratePartition(peer int) {
	if peer == -1 {
		peer = rand.Int() % 5
	}
	peer2 := (peer + rand.Int()%4 + 1) % 5
	fmt.Printf("Partition: [%d],[%d]\n", peer, peer2)
	block := make([]bool, len(Hosts))
	for i := len(Hosts) - 1; i >= 0; i-- {
		for j := len(Hosts) - 1; j >= 0; j-- {
			block[j] = shouldBlock(i, j, peer, peer2)
		}
		setBlock(i, block)
	}
}

func RecoverPartition() {
	block := make([]bool, len(Hosts))
	for i := len(Hosts) - 1; i >= 0; i-- {
		block[i] = false
	}
	for i := len(Hosts) - 1; i >= 0; i-- {
		setBlock(i, block)
	}
}

func GenerateCrash(peer int) {
	for {
		res := SendClientRequest(Hosts[peer]+"/client/crash", make([]byte, 0))
		if string(res) == "ok" {
			break
		}
	}
}
func RecoverCrash(peer int) {
	for {
		res := SendClientRequest(Hosts[peer]+"/client/recover", make([]byte, 0))
		if string(res) == "ok" {
			break
		}
	}
}

func GenerateDelay(peer int) int {
	if peer == -1 {
		peer = rand.Int() % 5
	}
	delay := make([]int, len(Hosts))
	for i := len(Hosts) - 1; i >= 0; i-- {
		if i != peer {
			delay[i] = 200 + rand.Int()%800
		} else {
			delay[i] = 0
		}
	}
	SendDelayRequest(peer, delay)
	return peer
}

func RecoverDelay(peer int) {
	delay := make([]int, len(Hosts))
	for i := len(Hosts) - 1; i >= 0; i-- {
		delay[i] = 0
	}
	SendDelayRequest(peer, delay)
}

func RandomOperator() int {
	switch rand.Int() % 4 {
	case 0:
		return global.SET
	case 1:
		return global.GET
	case 2:
		return global.DELETE
	default:
		return global.NULL
	}
}

func SendDelayRequest(peer int, delay []int) {
	var args global.DelayArgs
	args.Delay = delay
	body, _ := json.Marshal(&args)
	for {
		res := SendClientRequest(Hosts[peer]+"/client/recover", body)
		if string(res) == "ok" {
			break
		}
	}
}
