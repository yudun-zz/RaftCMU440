package raft

import (
	"github.com/yudun/RaftCMU440/rpc/raftrpc"
)

type raftNode struct {
	log []raftrpc.LogEntry
	// TODO: Implement this!
}

// Desc:
// NewRaftNode creates a new RaftNode. This function should return only when
// all nodes have joined the ring, and should return a non-nil error if this node
// could not be started in spite of dialing any other nodes numRetries times.
//
// Params:
// myHostPort: the hostport string of this new node. We use tcp in this project.
//			   	Note: Please listen to this port rather than hostMap[srvId]
// hostMap: a map from all node IDs to their hostports.
// numNodes: the number of nodes in the ring
// numRetries: if we can't connect with some nodes in hostMap after numRetries attempts,
// 				an error should be returned
// srvId: the id of this node
// heartBeatInterval: the Heart Beat Interval when this node becomes leader. In millisecond.
// electionTimeoutLow, electionTimeoutHigh: The election timeout for this node should be an
// random integer within range [electionTimeoutLow, electionTimeoutHigh]
func NewRaftNode(myHostPort string, hostMap map[int]string,
	numNodes, srvId, numRetries, heartBeatInterval,
	electionTimeoutLow, electionTimeoutHigh int) (RaftNode, error) {
	// TODO: Implement this!
	rn := &raftNode{}
	return rn, nil
}

// Desc:
// Propose initializes proposing a new operation, and replies with the
// result of committing this operation. Propose should not return until
// this operation has been committed, or this node is not leader now.
//
// If the we put a new <k, v> pair or deleted an existing <k, v> pair
// successfully, it should return OK; If it tries to delete an non-existing
// key, a KeyNotFound should be returned; If this node is not leader now,
// it should return WrongNode as well as the currentLeader id.
//
// Params:
// args: the operation to propose
// reply: as specified in Desc
func (rn *raftNode) Propose(args *raftrpc.ProposeArgs, reply *raftrpc.ProposeReply) error {
	// TODO: Implement this!
	return nil
}

// Desc:
// GetValue looks up the value for a key, and replies with the value or with
// the Status KeyNotFound.
//
// Params:
// args: the key to check
// reply: the value and status for this lookup of the given key
func (rn *raftNode) GetValue(args *raftrpc.GetValueArgs, reply *raftrpc.GetValueReply) error {
	// TODO: Implement this!
	return nil
}

// Desc:
// Receive a RecvRequestVote message from another Raft Node. Check the paper for more details.
//
// Params:
// args: the RequestVote Message, you must include From(src node id) and To(dst node id) when
// you call this API
// reply: the RequestVote Reply Message
func (rn *raftNode) RecvRequestVote(args *raftrpc.RequestVoteArgs, reply *raftrpc.RequestVoteReply) error {
	// TODO: Implement this!
	return nil
}

// Desc:
// Receive a RecvAppendEntries message from another Raft Node. Check the paper for more details.
//
// Params:
// args: the AppendEntries Message, you must include From(src node id) and To(dst node id) when
// you call this API
// reply: the AppendEntries Reply Message
func (rn *raftNode) RecvAppendEntries(args *raftrpc.AppendEntriesArgs, reply *raftrpc.AppendEntriesReply) error {
	// TODO: Implement this!
	return nil
}

// Desc:
// Set both the the electionTimeoutLow and electionTimeoutHigh of this node to be args.Timeout.
// You also need to stop current timer and reset it to fire after args.Timeout milliseconds.
//
// Please note that after this reset, ** if the electionTimer times out and restart and no new
// timeout duration re-generate is required (for example, if follower timeout and becomes
// candidate, we need to re-generate the timeout duration; but if the follower receives heartbeat
// from the leader, we don't generate new timeout duration -- you should restart using previous
// timeout duration) **, you should still set the timeout duration as args.Timeout milliseconds.
//
// Params:
// args: the election timeout duration
// reply: no use
func (rn *raftNode) RecvSetElectionTimeout(args *raftrpc.SetElectionTimeoutArgs, reply *raftrpc.SetElectionTimeoutReply) error {
	// TODO: Implement this!
	return nil
}

// Desc:
// Set heartBeatInterval as args.Interval milliseconds.
// You also need to stop current ticker and reset it to fire every args.Interval milliseconds.
//
// Params:
// args: the heartbeat duration
// reply: no use
func (rn *raftNode) RecvSetHeartBeatInterval(args *raftrpc.SetHeartBeatIntervalArgs, reply *raftrpc.SetHeartBeatIntervalReply) error {
	// TODO: Implement this!
	return nil
}
