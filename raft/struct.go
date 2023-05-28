package raft

import (
	"raft/global"
	"sync"
	"time"
)

type Raft struct {
	Mu          sync.Mutex
	Me          int
	State       int
	CurrentTerm int
	VotedFor    int
	Logs        []global.Log
	CommitIndex int
	LastApplied int
	NextIndex   []int
	MatchIndex  []int
	Votes       int
	Alive       bool

	TimeoutTimer   *time.Timer
	HeartBeatTimer []*time.Timer

	CheckMatchIndex  bool
	CheckLastApplied bool
	ApplyChan        chan global.ApplyMsg
}
