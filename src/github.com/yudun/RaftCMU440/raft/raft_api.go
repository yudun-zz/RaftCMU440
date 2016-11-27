package raft

import (
	"github.com/yudun/RaftCMU440/rpc/raftrpc"
)

type RaftNode interface {
	Propose(args *raftrpc.ProposeArgs, reply *raftrpc.ProposeReply) error
	GetValue(args *raftrpc.GetValueArgs, reply *raftrpc.GetValueReply) error
	RecvSetElectionTimeout(args *raftrpc.SetElectionTimeoutArgs, reply *raftrpc.SetElectionTimeoutReply) error
	RecvSetHeartBeatInterval(args *raftrpc.SetHeartBeatIntervalArgs, reply *raftrpc.SetHeartBeatIntervalReply) error
}
