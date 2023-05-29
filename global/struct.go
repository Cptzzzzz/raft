package global

type OperateArgs struct {
	Operator int    `json:"operator"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type OperateReply struct {
	Index  int  `json:"index"`
	Term   int  `json:"term"`
	Ok     bool `json:"ok"`
	Leader int  `json:"leader"`
}

type RequestVoteArgs struct {
	Term         int `json:"term"`
	CandidateId  int `json:"candidateId"`
	LastLogIndex int `json:"lastLogIndex"`
	LastLogTerm  int `json:"lastLogTerm"`
}

type RequestVoteReply struct {
	Term        int  `json:"term"`
	VoteGranted bool `json:"voteGranted"`
}

type AppendEntriesArgs struct {
	Term         int   `json:"term"`
	LeaderId     int   `json:"leaderId"`
	PrevLogIndex int   `json:"prevLogIndex"`
	PrevLogTerm  int   `json:"prevLogTerm"`
	Entries      []Log `json:"entries"`
	LeaderCommit int   `json:"leaderCommit"`
}

type AppendEntriesReply struct {
	Term          int  `json:"term"`
	Success       bool `json:"success"`
	ConflictTerm  int  `json:"conflictTerm"`
	ConflictIndex int  `json:"conflictIndex"`
}

type StateReply struct {
	State       string `json:"state"`
	CurrentTerm int    `json:"currentTerm"`
	VotedFor    int    `json:"votedFor"`
	Logs        []Log  `json:"logs"`
	CommitIndex int    `json:"commitIndex"`
	LastApplied int    `json:"lastApplied"`
	Alive       bool   `json:"alive"`
}

const (
	LEADER    = 1
	CANDIDATE = 2
	FOLLOWER  = 3
)

const (
	CREATE = iota
	UPDATE
	READ
	DELETE
	NULL
)

type Command struct {
	Operator int    `json:"operator"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Log struct {
	Index   int `json:"index"`
	Term    int `json:"term"`
	Command `json:"command"`
}

type ApplyMsg struct {
	Command
	CommandIndex int
}

type JudgeArgs struct {
	Command      `json:"command"`
	CommandIndex int `json:"commandIndex"`
	Peer         int `json:"peer"`
}
