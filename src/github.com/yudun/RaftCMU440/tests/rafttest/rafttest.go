package main

/*
*							A RAFT TESTING FRAMEWORK
*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*	This is a testing framework for Raft <https://raft.github.io>.
*
*	The basic idea of this test framework is to set up a proxy sever upon each raft node, and do some
*	bookkeeping in the proxy server to check if all the message are passed as expected. To ensure the
*	message are sent in expected order, we need to predefined the delay of each message, I store this
*	information in requestVoteSchema and appendEntriesSchema. Both of them are map from a message
*	description to an integer value indicating the delay of this message. The format of the message
*	description is:
*
*		t<term>:<node1><messagedirection><node2> #<order number of this in this term>
*
*	If this message is a response message, the <messagedirection> will be <-, otherwise, it is either
*	requestVote or appendEntries, the <messagedirection> will be ->. And if the value is negative, it
*	indicates we drop this message.
*
*	For example, in appendEntriesSchema, entry "t1:0->4 #1": 600 indiates the 1st appendEntries
*	message sent from node0 to node4 in term 1 will be delayed by 600ms; in requestVoteSchema, entry
*	"t1:0<-1 #2": -1 indiates the 2nd requestVote response sent from node 1 to node 0 in term 1 will
*	be dropped.
*
*	All undeclared messages will have no delay and will be be dropped.
*
*	To add tests to this framework, you need to:
*	1. construct the delay scenario in requestVoteSchema and appendEntriesSchema;
*	2. add expected events to the events array in each test;
*	3. broadcast this scenario to each proxy server by calling:
*			pt.RunAllRaftRound(requestVoteSchema, appendEntriesSchema, events, roundChan)
*	4. setup the desired election timeout and heartbeat interval for each node
*	5. using time.Sleep and pt.Propose to send out desired propose properly;
*	6. using checkGetValueAll, checkProposeRes and listen to the roundChan channel for the test
*	result from all nodes to check if there is something unexpceted.
*
*	Currently this testing framework contains 10 tests for leaderElection and 4 tests for log
*	replication.
*
*	Please feel free to contact me if you have any questions regarding to this framework
*
**/

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/yudun/RaftCMU440/rpc/raftproxyrpc"
	"github.com/yudun/RaftCMU440/tests/raftproxy"
)

type raftTester struct {
	myPort   int
	numNodes int
	cliMap   map[int]*rpc.Client
	px       []raftproxy.RaftProxy
}

type testFunc struct {
	name string
	f    func(chan bool)
}

var (
	numNodes   = flag.Int("N", 1, "number of raft nodes in the ring")
	proxyPorts = flag.String("proxyports", "", "proxy ports of all nodes")
	testRegex  = flag.String("t", "", "test to run")
	passCount  int
	failCount  int
	timeout    = 20
	pt         *raftTester
)

type ProposeResult struct {
	nodeId int
	reply  *raftproxyrpc.ProposeReply
	err    error
}

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)
var LOGF *os.File

func initRaftTester(allProxyPorts string, numNodes int) (*raftTester, error) {
	tester := new(raftTester)
	tester.numNodes = numNodes

	proxyPortList := strings.Split(allProxyPorts, ",")

	// Create RPC connection to raft nodes.
	cliMap := make(map[int]*rpc.Client)
	for i, p := range proxyPortList {
		cli, err := rpc.DialHTTP("tcp", "localhost:"+p)
		if err != nil {
			return nil, fmt.Errorf("could not connect to node %d", i)
		}
		cliMap[i] = cli
	}

	tester.cliMap = cliMap
	return tester, nil
}

func (pt *raftTester) GetValue(key string, nodeID int) (*raftproxyrpc.GetValueReply, error) {
	args := &raftproxyrpc.GetValueArgs{Key: key}
	var reply raftproxyrpc.GetValueReply
	err := pt.cliMap[nodeID].Call("RaftNode.GetValue", args, &reply)
	return &reply, err
}

func (pt *raftTester) ProposeAll(op raftproxyrpc.Operation, key string, value interface{},
	proposeRes chan ProposeResult) {
	for i := 0; i < pt.numNodes; i++ {
		go pt.Propose(op, key, value, i, proposeRes)
	}
}

func (pt *raftTester) Propose(op raftproxyrpc.Operation, key string, value interface{}, nodeID int,
	proposeRes chan ProposeResult) {
	args := &raftproxyrpc.ProposeArgs{Op: op, Key: key, V: value}
	var reply raftproxyrpc.ProposeReply
	err := pt.cliMap[nodeID].Call("RaftNode.Propose", args, &reply)
	proposeRes <- ProposeResult{
		nodeId: nodeID,
		reply:  &reply,
		err:    err,
	}
}

func (pt *raftTester) SetAllElectTimeout(timeouts []int) {
	for idx, timeout := range timeouts {
		go pt.SetElectTimeout(timeout, idx)
	}
}

func (pt *raftTester) SetElectTimeout(timeout, nodeID int) {
	args := &raftproxyrpc.SetElectionTimeoutArgs{
		Timeout: timeout,
	}
	var reply raftproxyrpc.SetElectionTimeoutReply
	pt.cliMap[nodeID].Call("RaftNode.RecvSetElectionTimeout", args, &reply)
}

func (pt *raftTester) SetHeartBeatInterval(timeout, nodeID int) {
	args := &raftproxyrpc.SetHeartBeatIntervalArgs{
		Interval: timeout,
	}
	var reply raftproxyrpc.SetHeartBeatIntervalReply
	pt.cliMap[nodeID].Call("RaftNode.RecvSetHeartBeatInterval", args, &reply)
}

/**
 * staff node-specific functions
 */
func (pt *raftTester) RunAllRaftRound(requestVoteSchema, appendEntriesSchema map[string]int,
	events []raftproxyrpc.Event, roundChan chan error) {
	for i := 0; i < pt.numNodes; i++ {
		go pt.RunRaftRound(requestVoteSchema, appendEntriesSchema, events, i, roundChan)
	}
}

func (pt *raftTester) RunRaftRound(requestVoteSchema, appendEntriesSchema map[string]int,
	events []raftproxyrpc.Event, nodeID int, roundChan chan error) {

	args := &raftproxyrpc.CheckEventsArgs{RequestVoteSchema: requestVoteSchema,
		AppendEntriesSchema: appendEntriesSchema,
		ExpectedEvents:      events}

	var reply raftproxyrpc.CheckEventsReply

	err := pt.cliMap[nodeID].Call("RaftNode.CheckEvents", args, &reply)
	if err != nil {
		roundChan <- err
		return
	} else {
		if reply.Success {
			roundChan <- nil
		} else {
			roundChan <- errors.New(reply.ErrMsg)
		}
	}
}

func checkProposeRes(res, expected ProposeResult, nodeId int) bool {
	if res.reply.Status == expected.reply.Status {
		if res.reply.Status == raftproxyrpc.WrongNode &&
			res.reply.CurrentLeader != expected.reply.CurrentLeader {
			LOGE.Printf("FAIL: incorrect current Leader from Propose to node %d. Got %v, expected %v.\n",
				nodeId, res.reply.CurrentLeader, expected.reply.CurrentLeader)
			return false
		}
	} else {
		LOGE.Printf("FAIL: incorrect Status from Propose to node %d. Got %v, expected %v.\n",
			nodeId, res.reply.Status, expected.reply.Status)
		return false
	}

	return true
}

func checkGetValueAll(key string, status []raftproxyrpc.Status, values []uint32) bool {
	for id, _ := range pt.cliMap {
		r, err := pt.GetValue(key, id)
		if err != nil {
			printFailErr("GetValue", err)
			return false
		}

		if r.Status != status[id] {
			LOGE.Printf("FAIL: Node %d: incorrect status from GetValue on key %s. Got %v, expected %v.\n",
				id, key, r.Status, status[id])
			return false
		} else {
			if r.Status == raftproxyrpc.KeyFound && r.V.(uint32) != values[id] {
				LOGE.Printf("FAIL: Node %d: incorrect value from GetValue on key %s. Got %d, expected %d.\n",
					id, key, r.V.(uint32), values[id])
				return false
			}
		}
	}
	return true
}

func runTest(t testFunc, doneChan chan bool) {
	go t.f(doneChan)

	var pass bool
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		LOGE.Println("FAIL: test timed out")
		pass = false
	case pass = <-doneChan:
	}

	var res string
	if pass {
		res = t.name + " PASS"
		fmt.Println(res)
		log.Println(res)
		passCount++
	} else {
		res = t.name + " FAIL"
		fmt.Println(res)
		log.Println(res)
		failCount++
	}
}

func printFailErr(fname string, err error) {
	LOGE.Printf("FAIL: error on %s - %s", fname, err)
}

func printFailHint(hint string) {
	LOGE.Printf("\nHint: %s\n", hint)
}

func main() {
	LOGF, err := os.OpenFile("rafttest.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Errorf("error opening file: %v", err)
	}
	defer LOGF.Close()

	log.SetOutput(LOGF)

	btests := []testFunc{
		// tests for leader election
		{"testOneCandidateOneRoundElection", testOneCandidateOneRoundElection},
		{"testOneCandidateStartTwoElection", testOneCandidateStartTwoElection},
		{"testTwoCandidateForElection", testTwoCandidateForElection},
		{"testSplitVote1", testSplitVote1},
		{"testSplitVote2", testSplitVote2},
		{"testAllForElection1", testAllForElection1},
		{"testAllForElection2", testAllForElection2},
		{"testLeaderRevertToFollower1", testLeaderRevertToFollower1},
		{"testLeaderRevertToFollower2", testLeaderRevertToFollower2},
		{"testLeaderRevertToFollower3", testLeaderRevertToFollower3},

		// tests for log replication
		{"testOneSimplePut", testOneSimplePut},
		{"testOneSimpleUpdate", testOneSimpleUpdate},
		{"testOneSimpleDelete", testOneSimpleDelete},
		{"testDeleteNonExistKey", testDeleteNonExistKey},
	}

	flag.Parse()

	// Run the tests with a single tester
	raftTester, err := initRaftTester(*proxyPorts, *numNodes)
	if err != nil {
		LOGE.Fatalln("Failed to initialize test:", err)
	}
	pt = raftTester

	doneChan := make(chan bool)
	for _, t := range btests {
		if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
			fmt.Printf("Running %s:\n", t.name)
			runTest(t, doneChan)
		}
	}
}
