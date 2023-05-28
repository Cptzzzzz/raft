package raft

import (
	"fmt"
	"log"
	"raft/global"
)

func DPrintf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func (rf *Raft) DPrintf(format string, a ...interface{}) {
	log.Printf(rf.stateMessage()+format, a...)
}

func (rf *Raft) stateMessage() string {
	return fmt.Sprintf("Term %d\t%s[%d]\t", rf.CurrentTerm, rf.stateString(), rf.Me)
}

func (rf *Raft) stateString() string {
	switch rf.State {
	case global.LEADER:
		return "Leader    "
	case global.CANDIDATE:
		return "Candidate "
	case global.FOLLOWER:
		return "Follower  "
	default:
		return ""
	}
}
