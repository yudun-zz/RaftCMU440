package main

/*
*	AUTHOR:
*			Shimin Wang <wyudun@gmail.com>
*
*	DESCRIPTION:
*			Runner for launching new raft instance
 */

import (
	"flag"
	"github.com/yudun/RaftCMU440/raft"
	"log"
	"strconv"
	"strings"
)

var (
	ports                      = flag.String("ports", "", "ports for all raft nodes")
	myport                     = flag.Int("myport", 9990, "port this raft node should listen to")
	numNodes                   = flag.Int("N", 1, "the number of nodes in the ring")
	nodeID                     = flag.Int("id", 0, "node ID must match index of this node's port in the ports list")
	numRetries                 = flag.Int("retries", 5, "number of times a node should retry dialing another node")
	heartBeatInterval          = flag.Int("heartBeatInterval", 10000, "heartBeat Interval in milliseconds")
	electionTimeoutIntervalLow = flag.Int("electionTimeoutIntervalLowerBound", 10000,
		"election Timeout LowerBound in milliseconds")
	electionTimeoutIntervalHigh = flag.Int("electionTimeoutIntervalHigherBound", 10000,
		"election Timeout HigherBound in milliseconds")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
}

func main() {
	flag.Parse()

	portStrings := strings.Split(*ports, ",")

	hostMap := make(map[int]string)
	for i, port := range portStrings {
		hostMap[i] = "localhost:" + port
	}

	// Create and start the Raft Node.
	_, err := raft.NewRaftNode("localhost:"+strconv.Itoa(*myport), hostMap, *numNodes,
		*nodeID, *numRetries, *heartBeatInterval, *electionTimeoutIntervalLow, *electionTimeoutIntervalHigh)

	if err != nil {
		log.Fatalln("Failed to create raft node:", err)
	}

	// Run the raft node forever.
	select {}
}
