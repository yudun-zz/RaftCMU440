package raftproxy

/*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*			Runner for launching new raft instance
 */

import (
	"fmt"
	"github.com/yudun/RaftCMU440/rpc/raftproxyrpc"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type RaftProxy interface {
	Reset(args *raftproxyrpc.CheckEventsArgs)
	Propose(args *raftproxyrpc.ProposeArgs, reply *raftproxyrpc.ProposeReply) error
	GetValue(args *raftproxyrpc.GetValueArgs, reply *raftproxyrpc.GetValueReply) error
	RecvRequestVote(args *raftproxyrpc.RequestVoteArgs, reply *raftproxyrpc.RequestVoteReply) error
	RecvAppendEntries(args *raftproxyrpc.AppendEntriesArgs, reply *raftproxyrpc.AppendEntriesReply) error
	CheckEvents(args *raftproxyrpc.CheckEventsArgs, reply *raftproxyrpc.CheckEventsReply) error
}

type logger struct {
	RequestVoteSchema, AppendEntriesSchema map[string]int
	expectedEvents                         []raftproxyrpc.Event
	events                                 []raftproxyrpc.Event
	newEvent                               chan raftproxyrpc.Event
}

type proxy struct {
	srv    *rpc.Client
	srvID  int
	logger *logger
}

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

var INFINITY_TIME = 100 * time.Minute

var keyFrenquencyMap = map[string]int{}
var keyFrenquencyMapLock = &sync.Mutex{}

/**
 * Proxy validates the requests going into a RaftNode and the responses coming out of it.  * It logs errors that occurs during a test.
 */
func NewProxy(nodePort, myPort int, srvId int) (RaftProxy, error) {
	p := new(proxy)
	p.logger = &logger{
		expectedEvents: []raftproxyrpc.Event{},
		events:         []raftproxyrpc.Event{},
		newEvent:       make(chan raftproxyrpc.Event, 1000),
	}
	p.srvID = srvId

	// Start server
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", myPort))
	if err != nil {
		LOGE.Println("Failed to listen:", err)
		return nil, err
	}

	// Create RPC connection to raft node.
	srv, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%d", nodePort))
	if err != nil {
		LOGE.Println("Failed to dial node %d", nodePort)
		return nil, err
	}
	p.srv = srv

	// register RPC
	rpc.RegisterName("RaftNode", raftproxyrpc.Wrap(p))
	rpc.HandleHTTP()
	go http.Serve(l, nil)

	fmt.Printf("node %d Proxy started\n", srvId)
	return p, nil
}

func (p *proxy) getExpectedEventsForThisNode(expectedEvents []raftproxyrpc.Event) []raftproxyrpc.Event {
	res := []raftproxyrpc.Event{}
	for _, event := range expectedEvents {
		if event.To == -1 {
			event.To = p.srvID
		}
		if event.From == -1 {
			event.From = p.srvID
		}

		if !event.IsResponse {
			if event.To == p.srvID && event.From != event.To {
				res = append(res, event)
			}
		} else {
			if event.From == p.srvID && event.From != event.To {
				res = append(res, event)
			}
		}
	}
	return res
}

func (p *proxy) Reset(args *raftproxyrpc.CheckEventsArgs) {
	p.logger = &logger{
		RequestVoteSchema:   args.RequestVoteSchema,
		AppendEntriesSchema: args.AppendEntriesSchema,
		expectedEvents:      p.getExpectedEventsForThisNode(args.ExpectedEvents),
		events:              []raftproxyrpc.Event{},
		newEvent:            make(chan raftproxyrpc.Event, 1000),
	}
}

func getEventsString(events []raftproxyrpc.Event) string {
	res := ""
	for _, event := range events {
		res += event.ToString() + "\n"
	}
	return res
}

func (l *logger) checkLogger() (bool, string) {
	sortedEvents := []string{}
	sortedExpectedEvents := []string{}

	for i := 0; i < len(l.events); i++ {
		sortedEvents = append(sortedEvents, l.events[i].ToString())
		sortedExpectedEvents = append(sortedExpectedEvents, l.expectedEvents[i].ToString())
	}

	sort.Strings(sortedEvents)
	sort.Strings(sortedExpectedEvents)

	if strings.Join(sortedEvents, "") != strings.Join(sortedExpectedEvents, "") {
		return false, fmt.Sprintf("\nExpected:\n%sBut get:\n%s",
			getEventsString(l.expectedEvents),
			getEventsString(l.events))
	}

	return true, ""
}

func (p *proxy) CheckEvents(args *raftproxyrpc.CheckEventsArgs, reply *raftproxyrpc.CheckEventsReply) error {
	p.Reset(args)
	newlogger := p.logger

	if len(newlogger.expectedEvents) == 0 {
		reply.Success = true
		return nil
	}

	for {
		select {
		case newEvent := <-newlogger.newEvent:

			// If you want to see the detail log trace for each test, please uncomment this line
			//fmt.Printf("%s\n",  newEvent.ToString())

			newlogger.events = append(newlogger.events, newEvent)
			if len(newlogger.events) == len(newlogger.expectedEvents) {
				reply.Success, reply.ErrMsg = newlogger.checkLogger()
				return nil
			}
		}
	}

	return nil
}

func (p *proxy) getDelayDuration(from, to, term int, isRequestVote, isReply bool) (bool, time.Duration) {
	newlogger := p.logger
	var key string
	var phase string
	if isRequestVote {
		phase = "rv "
	} else {
		phase = "ae "
	}

	if isReply {
		key = fmt.Sprintf("t%d:%d<-%d", term, to, from)
	} else {
		key = fmt.Sprintf("t%d:%d->%d", term, from, to)
	}

	phasekey := phase + key

	keyFrenquencyMapLock.Lock()
	if _, exist := keyFrenquencyMap[phasekey]; !exist {
		keyFrenquencyMap[phasekey] = 1
	} else {
		keyFrenquencyMap[phasekey]++
	}
	key = fmt.Sprintf("%s #%d", key, keyFrenquencyMap[phasekey])
	keyFrenquencyMapLock.Unlock()

	var delay int = 0
	if isRequestVote {
		if v, exist := newlogger.RequestVoteSchema[key]; exist {
			delay = v
		}
	} else {
		if v, exist := newlogger.AppendEntriesSchema[key]; exist {
			delay = v
		}
	}

	if delay < 0 {
		return true, time.Millisecond
	} else {
		return false, time.Duration(delay) * time.Millisecond
	}
}

func (p *proxy) Propose(args *raftproxyrpc.ProposeArgs, reply *raftproxyrpc.ProposeReply) error {
	err := p.srv.Call("RaftNode.Propose", args, reply)
	return err
}

func (p *proxy) GetValue(args *raftproxyrpc.GetValueArgs, reply *raftproxyrpc.GetValueReply) error {
	err := p.srv.Call("RaftNode.GetValue", args, reply)
	return err
}

func (p *proxy) RecvSetElectionTimeout(args *raftproxyrpc.SetElectionTimeoutArgs, reply *raftproxyrpc.SetElectionTimeoutReply) error {
	err := p.srv.Call("RaftNode.RecvSetElectionTimeout", args, reply)
	return err
}

func (p *proxy) RecvSetHeartBeatInterval(args *raftproxyrpc.SetHeartBeatIntervalArgs, reply *raftproxyrpc.SetHeartBeatIntervalReply) error {
	err := p.srv.Call("RaftNode.RecvSetHeartBeatInterval", args, reply)
	return err
}

func (p *proxy) RecvRequestVote(args *raftproxyrpc.RequestVoteArgs, reply *raftproxyrpc.RequestVoteReply) error {
	dropped, delay := p.getDelayDuration(args.From, args.To, args.Term, true, false)
	if dropped {
		fmt.Printf("node %d: dropped RequestVote from %d\n", args.To, args.From)
		time.Sleep(INFINITY_TIME)
		return nil
	}
	time.Sleep(delay)

	requestVoteEvent := raftproxyrpc.Event{
		From:         args.From,
		To:           args.To,
		Term:         args.Term,
		Msg:          raftproxyrpc.RequestVote,
		CandidateId:  args.From,
		LastLogIndex: args.LastLogIndex,
		LastLogTerm:  args.LastLogTerm,
		IsResponse:   false,
	}
	p.logger.newEvent <- requestVoteEvent

	err := p.srv.Call("RaftNode.RecvRequestVote", args, reply)

	dropped, delay = p.getDelayDuration(reply.From, reply.To, reply.Term, true, true)
	if dropped {
		fmt.Printf("node %d: dropped RequestVoteResponse to %d\n", reply.From, reply.To)
		time.Sleep(INFINITY_TIME)
		return nil
	}
	time.Sleep(delay)

	requestVoteResponseEvent := raftproxyrpc.Event{
		From:        reply.From,
		To:          reply.To,
		Term:        reply.Term,
		Msg:         raftproxyrpc.RequestVote,
		VoteGranted: reply.VoteGranted,
		IsResponse:  true,
	}

	p.logger.newEvent <- requestVoteResponseEvent
	return err
}

func (p *proxy) RecvAppendEntries(args *raftproxyrpc.AppendEntriesArgs, reply *raftproxyrpc.AppendEntriesReply) error {
	dropped, delay := p.getDelayDuration(args.From, args.To, args.Term, false, false)
	if dropped {
		fmt.Printf("node %d: dropped AppendEntries from %d\n", args.To, args.From)
		time.Sleep(INFINITY_TIME)
		return nil
	}
	time.Sleep(delay)

	appendEntriesEvent := raftproxyrpc.Event{
		From:         args.From,
		To:           args.To,
		Term:         args.Term,
		Msg:          raftproxyrpc.AppendEntries,
		LeaderId:     args.LeaderId,
		PrevLogIndex: args.PrevLogIndex,
		PrevLogTerm:  args.PrevLogTerm,
		Entries:      args.Entries,
		LeaderCommit: args.LeaderCommit,
		IsResponse:   false,
	}
	p.logger.newEvent <- appendEntriesEvent

	err := p.srv.Call("RaftNode.RecvAppendEntries", args, reply)

	dropped, delay = p.getDelayDuration(reply.From, reply.To, reply.Term, false, true)
	if dropped {
		fmt.Printf("node %d: dropped AppendEntriesResponse to %d\n", reply.From, reply.To)
		time.Sleep(INFINITY_TIME)
		return nil
	}
	time.Sleep(delay)

	appendEntriesResponseEvent := raftproxyrpc.Event{
		From:       reply.From,
		To:         reply.To,
		Term:       reply.Term,
		Msg:        raftproxyrpc.AppendEntries,
		Success:    reply.Success,
		MatchIndex: reply.MatchIndex,
		IsResponse: true,
	}

	p.logger.newEvent <- appendEntriesResponseEvent
	return err
}
