
A Raft Implementation Practice
==

This repository contains the starter code for a raft implementation and a suite of tests framework to help test your implementation. 

For what is Raft and how to implement the Raft Please refer to the annotated Raft paper in this repo.

In this project you will set up a **distributed replicated in-memory key-value store** using Raft to gurantee concensus. It support 2 operations:

1. `Put` will put a \<key, value\> pair into your store;
2. `Delete ` will try to delete a \<key, value\> pair from your store.



I made some simplication for this Raft implementation

1. You don't need to care about log compaction and snapshotting;
2. You can assume all nodes are robust and no node crash may happen — so you don't have to worry about membership changes;
3. No disk storage is considered in this implementation. You will have a `permenant` key-value store in memory. When you commit any log, you should directly apply that operation to this store, `committed` will be treated equally as` applied`. Consequently, the `lastApplied` argument mentioned in the paper unnecessary here.



Following instructions assume you have set your `GOPATH` to point to the repository's
root `RaftCMU440/` directory.

## Starter Code

The starter code for this project is organized roughly as follows:

```
src/github.com/yudun/RaftCMU440     
  raft/								implement Raft
  
  tests/							Source code for tests
  
  rpc/								Raft RPC helpers/constants
    
tests/                             	Shell scripts to run the tests
    rafttest.sh						Script for running all the tests
    rafttest_single.sh				Script for running any single test
```

## Instructions

### Compiling your code

To and compile your code, execute one or more of the following commands (the
resulting binaries will be located in the `$GOPATH/bin` directory):

```bash
go install github.com/yudun/RaftCMU440/runners/raftrunner
```

To simply check that your code compiles (i.e. without creating the binaries),
you can use the `go build` subcommand to compile an individual package as shown below:

```bash
# Build/compile the "raft" package.
go build path/to/raft

# A different way to build/compile the "raft" package.
go build github.com/yudun/RaftCMU440/raft
```

##### How to Write Go Code

If at any point you have any trouble with building, installing, or testing your code, the article
titled [How to Write Go Code](http://golang.org/doc/code.html) is a great resource for understanding
how Go workspaces are built and organized. You might also find the documentation for the
[`go` command](http://golang.org/cmd/go/) to be helpful. 



##### The `raftrunner` program

The `raftrunner` program creates and runs an instance of your
`RaftNode` implementation. Some example usage is provided below:

```bash
# Start a ring of three paxos nodes, where node 0 has port 9009, node 1 has port 9010, and so on.
./raftrunner -myport=9009 -ports=9009,9010,9011 -N=3 -id=0 -retries=5
./raftrunner -myport=9010 -ports=9009,9010,9011 -N=3 -id=1 -retries=5
./raftrunner -myport=9011 -ports=9009,9010,9011 -N=3 -id=2 -retries=5
```

You can read further descriptions of these flags by running 
```bash
./raftrunner -h
```

### Executing the tests

The tests for checkpoint are provided as bash shell scripts in the `tests` directory.
The scripts may be run from anywhere on your system (assuming your `GOPATH` has been set and
they are being executed on a 64-bit Mac OS X or Linux machine). Simply execute the following:

```bash
$GOPATH/tests/rafttest.sh
```

This will run all the tests (It may takes 200 ~ 300 seconds to finish)

If you only want to run a single test, please execute:

```bash
$GOPATH/tests/rafttest_single.sh <testname>
```

You can find the name of each test in  `github.com/yudun/RaftCMU440/tests/rafttest/rafttest.go` 

**Note:** You can uncommnet the log in func `CheckEvents` in `github.com/yudun/RaftCMU440/tests/raftproxy/raftproxy.go`  to see detailed log for each test.

