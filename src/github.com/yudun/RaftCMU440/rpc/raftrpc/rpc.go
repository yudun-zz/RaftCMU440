package raftrpc

/*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*			Signatures for the RPC in raft
 */

type RemoteRaftNode interface {
	// Called by servers using the raft node.
	Propose(args *ProposeArgs, reply *ProposeReply) error
	GetValue(args *GetValueArgs, reply *GetValueReply) error
	RecvSetElectionTimeout(args *SetElectionTimeoutArgs, reply *SetElectionTimeoutReply) error
	RecvSetHeartBeatInterval(args *SetHeartBeatIntervalArgs, reply *SetHeartBeatIntervalReply) error

	// Called by other Raft Nodes.
	RecvRequestVote(args *RequestVoteArgs, reply *RequestVoteReply) error
	RecvAppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) error
}

type RaftNode struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemoteRaftNode
}

// Wrap wraps t in a type-safe wrapper struct to ensure that only the desired
// methods are exported to receive RPCs.
func Wrap(t RemoteRaftNode) RemoteRaftNode {
	return &RaftNode{t}
}
