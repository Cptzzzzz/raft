package raft

import (
	"encoding/json"
	"math/rand"
	"raft/global"
	"time"
)

func Make(me int) *Raft {
	rf := &Raft{}
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	rf.Alive = true
	rf.Me = me
	rf.CurrentTerm = 0
	rf.State = global.FOLLOWER
	rf.VotedFor = -1
	rf.Votes = 0
	rf.Logs = make([]global.Log, 0)
	rf.Logs = append(rf.Logs, global.Log{
		Index: 0,
		Term:  0,
	})
	rf.HeartBeatTimer = make([]*time.Timer, len(global.Peers))
	rf.TimeoutTimer = time.AfterFunc(rf.electionTimeout(), rf.startElection)
	rf.ApplyChan = make(chan global.ApplyMsg)
	go rf.applier(rf.ApplyChan)
	return rf
}

func (rf *Raft) sendRequestVote(peer int, args global.RequestVoteArgs) {
	var reply global.RequestVoteReply
	body, _ := json.Marshal(&args)
	body = SendRequest(global.Peers[peer]+"/raft/request-vote", body)
	if body == nil {
		return
	}
	json.Unmarshal(body, &reply)
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	if !(rf.State == global.CANDIDATE &&
		rf.CurrentTerm == args.Term &&
		rf.Me == args.CandidateId) {
		return
	}
	if rf.CurrentTerm < reply.Term {
		//todo become follower
		rf.becomeFollower(reply.Term, false)
		return
	}
	if reply.VoteGranted {
		rf.Votes++
	}
	if rf.Votes <= len(global.Peers)/2 {
		return
	}
	rf.DPrintf("become leader")
	rf.Votes = 0
	rf.State = global.LEADER
	rf.stopTimer()
	rf.NextIndex = make([]int, len(global.Peers))
	rf.MatchIndex = make([]int, len(global.Peers))
	for i := len(global.Peers) - 1; i >= 0; i-- {
		rf.NextIndex[i] = len(rf.Logs)
		rf.MatchIndex[i] = 0
		if rf.Me != i {
			go rf.startHeartbeat(i, rf.CurrentTerm)
		}
	}
}

func (rf *Raft) startHeartbeat(peer, term int) {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	if rf.State != global.LEADER || rf.CurrentTerm != term {
		return
	}
	args := global.AppendEntriesArgs{
		Term:         rf.CurrentTerm,
		LeaderId:     rf.Me,
		PrevLogIndex: rf.NextIndex[peer] - 1,
		PrevLogTerm:  rf.Logs[rf.NextIndex[peer]-1].Term,
		Entries:      rf.Logs[rf.NextIndex[peer]:],
		LeaderCommit: rf.CommitIndex,
	}
	if len(args.Entries) > 0 && rf.MatchIndex[peer]+1 != rf.NextIndex[peer] {
		args.Entries = args.Entries[:1]
	}
	go rf.sendAppendEntries(peer, term, args)
	rf.resetHeartbeatTimer(peer, term)
}

func (rf *Raft) sendAppendEntries(peer, term int, args global.AppendEntriesArgs) {
	var reply global.AppendEntriesReply
	body, _ := json.Marshal(&args)
	body = SendRequest(global.Peers[peer]+"/raft/append-entries", body)
	if body == nil {
		return
	}
	json.Unmarshal(body, &reply)
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	if rf.State != global.LEADER || rf.CurrentTerm != term {
		return
	}
	if term < reply.Term {
		//todo become follower
		rf.becomeFollower(reply.Term, true)
		return
	}
	if rf.NextIndex[peer]-1 != args.PrevLogIndex {
		return
	}
	if reply.Success {
		rf.MatchIndex[peer] = args.PrevLogIndex + len(args.Entries)
		rf.NextIndex[peer] = rf.MatchIndex[peer] + 1
		if !rf.CheckMatchIndex {
			go rf.checkMatchIndex()
			rf.CheckMatchIndex = true
		}
	} else {
		flag := true
		if reply.ConflictTerm != -1 {
			for index := len(rf.Logs) - 1; index >= 0; index-- {
				if rf.Logs[index].Term == reply.ConflictTerm {
					flag = false
				}
				if rf.Logs[index].Term <= reply.ConflictTerm {
					break
				}
				rf.NextIndex[peer] = index
			}
		}
		if flag {
			rf.NextIndex[peer] = reply.ConflictIndex
		}
	}
	if rf.NextIndex[peer] != len(rf.Logs) || rf.MatchIndex[peer] != len(rf.Logs)-1 {
		go rf.startHeartbeat(peer, term)
		rf.stopHeartbeatTimer(peer)
	}
}
func (rf *Raft) startElection() {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	rf.DPrintf("start a election")
	rf.State = global.CANDIDATE
	rf.CurrentTerm++
	rf.VotedFor = rf.Me
	rf.Votes = 1
	args := global.RequestVoteArgs{
		Term:         rf.CurrentTerm,
		CandidateId:  rf.Me,
		LastLogIndex: len(rf.Logs) - 1,
		LastLogTerm:  rf.Logs[len(rf.Logs)-1].Term,
	}
	for i := len(global.Peers) - 1; i >= 0; i-- {
		if i != rf.Me {
			go rf.sendRequestVote(i, args)
		}
	}
	rf.resetTimer()
}

func (rf *Raft) electionTimeout() time.Duration {
	return time.Millisecond * time.Duration(400+rand.Intn(300))
}

func (rf *Raft) resetTimer() {
	rf.stopTimer()
	if rf.Alive {
		rf.TimeoutTimer = time.AfterFunc(rf.electionTimeout(), rf.startElection)
	}
}
func (rf *Raft) stopTimer() {
	if rf.TimeoutTimer != nil {
		rf.TimeoutTimer.Stop()
	}
}

func (rf *Raft) heartbeatTimeout() time.Duration {
	return time.Millisecond * time.Duration(100)
}

func (rf *Raft) resetHeartbeatTimer(peer, term int) {
	rf.stopHeartbeatTimer(peer)
	if rf.Alive {
		rf.HeartBeatTimer[peer] = time.AfterFunc(rf.heartbeatTimeout(), func() {
			rf.startHeartbeat(peer, term)
		})
	}
}

func (rf *Raft) stopHeartbeatTimer(peer int) {
	if rf.HeartBeatTimer[peer] != nil {
		rf.HeartBeatTimer[peer].Stop()
	}
}

func (rf *Raft) becomeFollower(term int, reset bool) {
	if rf.State == global.LEADER {
		for i := len(global.Peers) - 1; i >= 0; i-- {
			if rf.Me != i {
				rf.stopHeartbeatTimer(i)
			}
		}
	}
	rf.CurrentTerm = term
	if rf.State == global.LEADER || reset {
		rf.resetTimer()
	}
	rf.State = global.FOLLOWER
	rf.VotedFor = -1
	rf.Votes = 0
}
func (rf *Raft) checkMatchIndex() {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	rf.CheckMatchIndex = false
	if rf.State != global.LEADER {
		return
	}
	last := rf.CommitIndex
	for N := len(rf.Logs) - 1; N >= rf.CommitIndex+1; N++ {
		votes := 0
		for i := len(global.Peers) - 1; i >= 0; i-- {
			if rf.MatchIndex[i] >= N || i == rf.Me {
				votes++
			}
		}
		if votes > len(global.Peers)/2 && rf.Logs[N].Term == rf.CurrentTerm {
			rf.CommitIndex = N
			break
		}
	}
	if last != rf.CommitIndex && !rf.CheckLastApplied {
		go rf.checkLastApplied()
		rf.CheckLastApplied = true
	}
}

func (rf *Raft) checkLastApplied() {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	rf.CheckLastApplied = false
	if rf.CurrentTerm != rf.Logs[rf.CommitIndex].Term {
		return
	}
	for rf.CommitIndex > rf.LastApplied {
		rf.LastApplied++
		rf.DPrintf("ready to send ApplyMsg")
		rf.ApplyChan <- global.ApplyMsg{
			Command:      rf.Logs[rf.LastApplied].Command,
			CommandIndex: rf.LastApplied,
		}
		rf.DPrintf("sent ApplyMsg")
	}
}

func (rf *Raft) AppendEntries(args *global.AppendEntriesArgs, reply *global.AppendEntriesReply) {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	reply.Success = false
	if args.Term < rf.CurrentTerm {
		reply.Term = rf.CurrentTerm
		return
	}
	rf.becomeFollower(args.Term, true)
	rf.VotedFor = args.LeaderId
	reply.Term = rf.CurrentTerm
	if len(rf.Logs) <= args.PrevLogIndex {
		reply.ConflictIndex = len(rf.Logs) - 1
		reply.ConflictTerm = -1
		return
	}
	if args.PrevLogTerm != rf.Logs[args.PrevLogIndex].Term {
		reply.ConflictTerm = rf.Logs[args.PrevLogIndex].Term
		for index := args.PrevLogIndex; index >= 0 && rf.Logs[index].Term == reply.ConflictTerm; index-- {
			if rf.Logs[index].Term == reply.ConflictTerm {
				reply.ConflictIndex = index
			}
		}
		rf.Logs = rf.Logs[:args.PrevLogIndex]
		return
	}
	for index, val := range args.Entries {
		if len(rf.Logs) == index+1+args.PrevLogIndex ||
			val.Term != rf.Logs[index+1+args.PrevLogIndex].Term {
			rf.Logs = rf.Logs[:index+1+args.PrevLogIndex]
			rf.Logs = append(rf.Logs, args.Entries[index:]...)
			break
		}
	}
	if args.LeaderCommit > rf.CommitIndex {
		rf.CommitIndex = min(args.LeaderCommit, len(rf.Logs)-1)
		if !rf.CheckLastApplied {
			rf.CheckLastApplied = true
			go rf.checkLastApplied()
		}
	}
	reply.Success = true
}

func (rf *Raft) RequestVote(args *global.RequestVoteArgs, reply *global.RequestVoteReply) {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	reply.VoteGranted = false
	if rf.CurrentTerm > args.Term {
		reply.Term = rf.CurrentTerm
		return
	}
	if rf.CurrentTerm < args.Term {
		rf.becomeFollower(args.Term, false)
	}
	reply.Term = rf.CurrentTerm
	if rf.State == global.FOLLOWER && (rf.VotedFor == -1 || rf.VotedFor == args.CandidateId) {
		if args.LastLogTerm > rf.Logs[len(rf.Logs)-1].Term ||
			(args.LastLogTerm == rf.Logs[len(rf.Logs)-1].Term && args.LastLogIndex >= len(rf.Logs)-1) {
			reply.VoteGranted = true
			rf.VotedFor = args.CandidateId
			rf.resetTimer()
		}
	}
}

func (rf *Raft) applier(ch chan global.ApplyMsg) {
	for {
		msg := <-ch
		rf.DPrintf("apply %v at %d", msg.Command, msg.CommandIndex)
		args := global.JudgeArgs{
			Command:      msg.Command,
			CommandIndex: msg.CommandIndex,
			Peer:         rf.Me,
		}
		body, _ := json.Marshal(&args)
		SendRequest(global.JudgeHost+"/msg", body)
	}
}

func (rf *Raft) Append(command global.Command) (int, int, bool, int) {
	rf.Mu.Lock()
	defer rf.Mu.Unlock()
	if rf.State != global.LEADER {
		return -1, rf.CurrentTerm, false, rf.VotedFor
	}
	entryLog := global.Log{
		Index:   len(rf.Logs),
		Term:    rf.CurrentTerm,
		Command: command,
	}
	rf.Logs = append(rf.Logs, entryLog)
	rf.MatchIndex[rf.Me] = entryLog.Index
	return entryLog.Index, rf.CurrentTerm, true, rf.Me
}

var Rf *Raft
