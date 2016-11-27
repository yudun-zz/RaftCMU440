package raftrpc

/*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*			Header files for the RPC in raft
 */

import "fmt"

// Status represents the status of a RPC's reply.
type Role int
type Operation int
type Status int

const (
	Follower Role = iota + 1
	Candidate
	Leader
)

const (
	Put Operation = iota + 1
	Delete
)

const (
	OK Status = iota + 1
	KeyFound
	KeyNotFound
	WrongNode
)

type ProposeArgs struct {
	Op  Operation
	Key string
	V   interface{} // Value for the Key if Put
}

type ProposeReply struct {
	CurrentLeader int
	Status        Status
}

type GetValueArgs struct {
	Key string
}

type GetValueReply struct {
	V      interface{}
	Status Status
}

type LogEntry struct {
	Term  int
	Op    Operation
	Key   string
	Value interface{}
}

func (l LogEntry) ToString() string {
	return fmt.Sprintf("%d:<%s,%v>", l.Term, l.Key, l.Value)
}

type RequestVoteArgs struct {
	From         int
	To           int
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

func (rv *RequestVoteArgs) ToString() string {
	return fmt.Sprintf("RequestVoteArgs: t%d:%d->%d, CandidateId:%d, LastLogIndex:%d, LastLogTerm:%d\n",
		rv.Term, rv.From, rv.To, rv.CandidateId, rv.LastLogIndex, rv.LastLogTerm)
}

type RequestVoteReply struct {
	From        int
	To          int
	Term        int
	VoteGranted bool
}

type AppendEntriesArgs struct {
	From         int
	To           int
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int
}

type AppendEntriesReply struct {
	From       int
	To         int
	Term       int
	Success    bool
	MatchIndex int
}

type SetElectionTimeoutArgs struct {
	Timeout int
}

type SetElectionTimeoutReply struct {
}

type SetHeartBeatIntervalArgs struct {
	Interval int
}

type SetHeartBeatIntervalReply struct {
}
