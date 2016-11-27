package main

/*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*			Tests for leader election
 */

import (
	"github.com/yudun/RaftCMU440/rpc/raftproxyrpc"
	"time"
)

/**
 * Tests only one candidate participant in the election with 1 simple round
 */
func testOneCandidateOneRoundElection(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0<-3 #1": -1,
		"t1:0<-4 #1": -1,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 1 vote for 0
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		// 2 vote for 0
		{Term: 1, From: 2, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		// 3, 4 vote for 0 dropped
		// 0 appendEntry for all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 1, From: -1, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 2000, 2000, 2000, 2000})

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testOneCandidateOneRoundElection", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Tests only one candidate participant in the election with 2 rounds
 */
func testOneCandidateStartTwoElection(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0->2 #1": -1,
		"t1:0->3 #1": -1,
		"t1:0->4 #1": -1,
		"t2:0<-1 #1": -1,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to 1
		{Term: 1, From: 0, To: 1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// requestVote to 2,3,4 lost
		// 1 vote for 0
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// after 1 seconds, 0 start new election with term 2
		// 0 requestVote to all
		{Term: 2, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 2, 3, 4 vote for 0
		{Term: 2, From: 2, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		{Term: 2, From: 3, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		{Term: 2, From: 4, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		// 1 vote for 0 dropped

		// 0 appendEntry for all
		{Term: 2, From: 0, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 2, From: -1, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 3000, 3000, 3000, 3000})

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testOneCandidateStartTwoElection", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Tests two candidates run for the election
 */
func testTwoCandidateForElection(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0->2 #1": -1,
		"t1:0->3 #1": -1,
		"t1:0->4 #1": -1,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to 1
		{Term: 1, From: 0, To: 1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 1 vote for 0
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// after 1 seconds, 2 start election with term 1
		// 2 requestVote to all
		{Term: 1, From: 2, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 0, 1 rej 2, 3, 4 vote for 2
		{Term: 1, From: 0, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 1, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 3, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		{Term: 1, From: 4, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 2 appendEntry for all
		{Term: 1, From: 2, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 2, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 1, From: -1, To: 2, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{2000, 3000, 3000, 4000, 4000})

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testTwoCandidateForElection", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Tests split vote for 2 candidates and election won by another candidate
 */
func testSplitVote1(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0->3 #1": 500,
		"t1:0->4 #1": 500,
		"t1:0<-2 #1": -1,
		"t1:3->0 #1": 500,
		"t1:3->1 #1": 500,
		"t1:3->2 #1": 500,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 1 vote for 0
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		// 2 dropped; 3,4 rej 0
		{Term: 1, From: 3, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 4, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},

		// 3 requestVote to all
		{Term: 1, From: 3, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 3, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 0, 1, 2 rej 3; 4 vote for 3
		{Term: 1, From: 0, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 1, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 2, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 4, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 2 requestVote to all with term 2
		{Term: 2, From: 2, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// all grant
		{Term: 2, From: -1, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 2 appendEntry for all
		{Term: 2, From: 2, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 2, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 2, From: -1, To: 2, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{2000, 3000, 3000, 2000, 3000})
	time.Sleep(2600 * time.Millisecond)
	go pt.SetElectTimeout(500, 2)

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testSplitVote1", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Tests split vote for 2 candidates and election won by one of the candidate
 */
func testSplitVote2(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0->3 #1": 600,
		"t1:0->4 #1": 600,
		"t1:0<-2 #1": -1,
		"t1:3->0 #1": 600,
		"t1:3->1 #1": 600,
		"t1:3->2 #1": 600,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 1 vote for 0
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		// 2 dropped; 3,4 rej 0
		{Term: 2, From: 3, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: 4, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},

		// 3 requestVote to all
		{Term: 1, From: 3, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 3, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 0, 1, 2 rej 3; 4 vote for 3
		{Term: 2, From: 0, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: 1, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: 2, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: 4, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 3 restart election, requestVote to all with term 2
		{Term: 2, From: 3, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 3, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// all grant
		{Term: 2, From: -1, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 3 appendEntry for all
		{Term: 2, From: 3, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 3, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 2, From: -1, To: 3, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{2000, 3000, 3000, 2000, 3000})
	time.Sleep(2100 * time.Millisecond)
	go pt.SetElectTimeout(300, 3)

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testSplitVote2", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Tests all nodes run for election and none gain majority in 1st round
 */
func testAllForElection1(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0->1 #1": 500, "t1:0->2 #1": 500, "t1:0->3 #1": 500, "t1:0->4 #1": 500,
		"t1:1->0 #1": 500, "t1:1->2 #1": 500, "t1:1->3 #1": 500, "t1:1->4 #1": 500,
		"t1:2->0 #1": 500, "t1:2->1 #1": 500, "t1:2->3 #1": 500, "t1:2->4 #1": 500,
		"t1:3->0 #1": 500, "t1:3->1 #1": 500, "t1:3->2 #1": 500, "t1:3->4 #1": 500,
		"t1:4->0 #1": 500, "t1:4->1 #1": 500, "t1:4->2 #1": 500, "t1:4->3 #1": 500,
		"t2:0->1 #1": 500, "t2:0->2 #1": 500, "t2:0->3 #1": 500, "t2:0->4 #1": 500,
		"t2:1->0 #1": 500, "t2:1->2 #1": 500, "t2:1->3 #1": 500, "t2:1->4 #1": 500,
		"t2:2->0 #1": 500, "t2:2->1 #1": 500, "t2:2->3 #1": 500, "t2:2->4 #1": 500,
		"t2:3->0 #1": 500, "t2:3->1 #1": 500, "t2:3->2 #1": 500, "t2:3->4 #1": 500,
		"t2:4->0 #1": 500, "t2:4->1 #1": 500, "t2:4->2 #1": 500, "t2:4->3 #1": 500,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// all requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 1, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 1, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 2, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 3, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 3, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 4, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},

		// all rej all
		{Term: 1, From: -1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 1, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},

		// all requestVote to all again
		{Term: 2, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 2, From: 1, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 1, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 2, From: 2, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 2, From: 3, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 3, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 2, From: 4, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},

		// all rej all again
		{Term: 2, From: -1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: -1, To: 1, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: -1, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: -1, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 2, From: -1, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 1000, 1000, 1000, 1000})

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testAllForElection1", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Tests all nodes run for election, one of them get elected in 2nd round
 */
func testAllForElection2(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0->1 #1": 500, "t1:0->2 #1": 500, "t1:0->3 #1": 500, "t1:0->4 #1": 500,
		"t1:1->0 #1": 500, "t1:1->2 #1": 500, "t1:1->3 #1": 500, "t1:1->4 #1": 500,
		"t1:2->0 #1": 500, "t1:2->1 #1": 500, "t1:2->3 #1": 500, "t1:2->4 #1": 500,
		"t1:3->0 #1": 500, "t1:3->1 #1": 500, "t1:3->2 #1": 500, "t1:3->4 #1": 500,
		"t1:4->0 #1": 500, "t1:4->1 #1": 500, "t1:4->2 #1": 500, "t1:4->3 #1": 500,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// all requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 1, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 1, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 2, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 2, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 3, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 3, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 1, From: 4, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},

		// all rej all
		{Term: 1, From: -1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 1, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 2, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 3, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},
		{Term: 1, From: -1, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: false, IsResponse: true},

		// after 600ms, 4 requestVote to All again
		{Term: 2, From: 4, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// all grant
		{Term: 2, From: -1, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 4 appendEntry for all
		{Term: 2, From: 4, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 4, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 2, From: -1, To: 4, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 1000, 1000, 1000, 1000})
	time.Sleep(1600 * time.Millisecond)
	go pt.SetElectTimeout(300, 4)

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testAllForElection2", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Node 0 becomes leader, and then convert to follower because node 4 requestVote with
 * a higher proposal number to 0
 */
func testLeaderRevertToFollower1(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t1:0<-4 #1": -1,
	}
	appendEntriesSchema := map[string]int{
		"t1:0->4 #1": -1,
	}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},

		// 1,2,3 grant 0; 4 dropped
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		{Term: 1, From: 2, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		{Term: 1, From: 3, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 0 appendEntry for 1,2,3; to 4 dropped
		{Term: 1, From: 0, To: 1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		{Term: 1, From: 0, To: 2, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		{Term: 1, From: 0, To: 3, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// 1,2,3 reply success
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
		{Term: 1, From: 2, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
		{Term: 1, From: 3, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},

		// after 400ms, 4 requestVote to all
		{Term: 2, From: 4, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// all grant
		{Term: 2, From: -1, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 4 appendEntries to all
		{Term: 2, From: 4, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 4, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success to 4
		{Term: 2, From: -1, To: 4, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 2000, 2000, 2000, 2000})
	time.Sleep(1100 * time.Millisecond)
	go pt.SetElectTimeout(300, 4)

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testLeaderRevertToFollower1", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Node 0 becomes leader, and then convert to follower because it receive higher
 * proposal number from the appendEntries reply from node 4
 */
func testLeaderRevertToFollower2(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t2:4->0 #1": -1,
		"t2:4->1 #1": -1,
		"t2:4->2 #1": -1,
	}
	appendEntriesSchema := map[string]int{
		"t1:0->4 #1": 600,
	}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// all grant 0
		{Term: 1, From: -1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 0 appendEntry for all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// 1,2,3 reply success
		{Term: 1, From: 1, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
		{Term: 1, From: 2, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
		{Term: 1, From: 3, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
		// 4 reply rej with term 2
		{Term: 2, From: 4, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: false, MatchIndex: 0, IsResponse: true},

		// after 400ms, 4 requestVote to 3; requestVote to 0,1,2 be dropped
		{Term: 2, From: 4, To: 3, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 3 grant
		{Term: 2, From: 3, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 2000, 2000, 2000, 2000})
	time.Sleep(1100 * time.Millisecond)
	go pt.SetElectTimeout(300, 4)

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testLeaderRevertToFollower2", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}

/**
 * Node 0 becomes leader, and then convert to follower because it receive higher
 * proposal number from the appendEntries request from new leader node 4
 */
func testLeaderRevertToFollower3(doneChan chan bool) {
	// following is the delayed schema for each message for RequestVote and AppendEntries
	requestVoteSchema := map[string]int{
		"t2:4->0 #1": -1,
		"t2:4->2 #1": -1,
	}
	appendEntriesSchema := map[string]int{}

	// following is the expected global events stream for this test
	events := []raftproxyrpc.Event{
		// 0 requestVote to all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 0, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// all grant 0
		{Term: 1, From: -1, To: 0, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 0 appendEntry for all
		{Term: 1, From: 0, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 0, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 1, From: -1, To: 0, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},

		// after 400ms, 4 requestVote to 1,3; requestVote to 0,2 be dropped
		{Term: 2, From: 4, To: 1, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		{Term: 2, From: 4, To: 3, Msg: raftproxyrpc.RequestVote,
			CandidateId: 4, LastLogIndex: 0, LastLogTerm: 0, IsResponse: false},
		// 1, 3 grant
		{Term: 2, From: 1, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},
		{Term: 2, From: 3, To: 4, Msg: raftproxyrpc.RequestVote,
			VoteGranted: true, IsResponse: true},

		// 4 appendEntry for all
		{Term: 2, From: 4, To: -1, Msg: raftproxyrpc.AppendEntries,
			LeaderId: 4, PrevLogIndex: 0, PrevLogTerm: 0, Entries: []raftproxyrpc.LogEntry{},
			LeaderCommit: 0, IsResponse: false},
		// all reply success
		{Term: 2, From: -1, To: 4, Msg: raftproxyrpc.AppendEntries,
			Success: true, MatchIndex: 0, IsResponse: true},
	}

	roundChan := make(chan error)
	pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
	time.Sleep(2 * time.Second)

	pt.SetAllElectTimeout([]int{1000, 2000, 2000, 2000, 2000})
	time.Sleep(1100 * time.Millisecond)
	go pt.SetElectTimeout(300, 4)

	for i := 0; i < *numNodes; i++ {
		err := <-roundChan
		if err != nil {
			printFailErr("testLeaderRevertToFollower3", err)
			doneChan <- false
			return
		}
	}

	doneChan <- true
}
