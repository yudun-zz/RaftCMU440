package raftproxyrpc

/*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*			Header files for the RPC in raft proxy
 */
import (
	"fmt"
	"strings"
)

// Status represents the status of a RPC's reply.
type Role int
type Operation int
type Status int
type MsgName int

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

const (
	RequestVote MsgName = iota + 1
	AppendEntries
)

type Event struct {
	Msg  MsgName
	From int
	To   int

	Term int

	// RequestVote
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int

	VoteGranted bool

	// AppendEntries
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry
	LeaderCommit int

	Success    bool
	MatchIndex int

	IsResponse bool
}

func (e Event) ToString() string {
	var res = ""
	if e.IsResponse {
		res += fmt.Sprintf("node %d -> %d: ", e.From, e.To)
	} else {
		res += fmt.Sprintf("node %d <- %d: ", e.To, e.From)
	}
	switch e.Msg {
	case RequestVote:
		if e.IsResponse {
			if e.VoteGranted {
				res += fmt.Sprintf("granted, term: %d", e.Term)
			} else {
				res += fmt.Sprintf("reject, term: %d", e.Term)
			}
		} else {
			res += fmt.Sprintf("RequestVote -- term: %d, candidateId: %d, lastLogIdx: %d, lastLogTerm: %d",
				e.Term, e.CandidateId, e.LastLogIndex, e.LastLogTerm)
		}
	case AppendEntries:
		if e.IsResponse {
			if e.Success {
				res += fmt.Sprintf("success, term: %d, matchIdx: %d", e.Term, e.MatchIndex)
			} else {
				res += fmt.Sprintf("fail, term: %d, matchIdx: %d", e.Term, e.MatchIndex)
			}
		} else {
			res += fmt.Sprintf("AppendEntries -- term: %d, leaderId: %d, prevLogIdx: %d, prevLogTerm: %d, "+
				"entries: %s, leaderCommit: %d", e.Term, e.LeaderId, e.PrevLogIndex, e.PrevLogTerm,
				EntriesToString(e.Entries), e.LeaderCommit)
		}
	}
	return res
}

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

func EntriesToString(l []LogEntry) string {
	s := []string{}
	for _, le := range l {
		s = append(s, le.ToString())
	}
	return fmt.Sprintf("[%s]", strings.Join(s, ","))
}

func (l LogEntry) ToString() string {
	var op string
	switch l.Op {
	case Put:
		op = "put"
	case Delete:
		op = "del"
	}
	return fmt.Sprintf("t%d:%s<%s,%v>", l.Term, op, l.Key, l.Value)
}

type RequestVoteArgs struct {
	From         int
	To           int
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
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

type CheckEventsArgs struct {
	RequestVoteSchema, AppendEntriesSchema map[string]int
	ExpectedEvents                         []Event
}

type CheckEventsReply struct {
	Success bool
	ErrMsg  string
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
