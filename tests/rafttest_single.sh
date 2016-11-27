#!/usr/bin/env bash

#!/bin/bash

if [ -z $GOPATH ]; then
    echo "FAIL: GOPATH environment variable is not set"
    exit 1
fi

if [ -n "$(go version | grep 'darwin/amd64')" ]; then
    GOOS="darwin_amd64"
elif [ -n "$(go version | grep 'linux/amd64')" ]; then
    GOOS="linux_amd64"
else
    echo "FAIL: only 64-bit Mac OS X and Linux operating systems are supported"
    exit 1
fi

# Build the student's raft node implementation.
# Exit immediately if there was a compile-time error.
go install github.com/yudun/RaftCMU440/runners/raftrunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Build the proxy runner
# Exit immediately if there was a compile-time error.
go install github.com/yudun/RaftCMU440/runners/raftproxyrunner
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

# Build the test binary to use to test the student's raft node implementation.
# Exit immediately if there was a compile-time error.
go install github.com/yudun/RaftCMU440/tests/rafttest
if [ $? -ne 0 ]; then
   echo "FAIL: code does not compile"
   exit $?
fi

if [ $# -le 1 ]; then
    # Pick distinct random ports between [10000, 20000).
    while ((i<11))
    do
       N=$(((RANDOM % 10000) + 10000))
       echo "${A[*]}" | grep $N && continue # if number already in the array
       A[$i]=$N
       ((i++))
    done

    NODE_PORT0=${A[0]}
    NODE_PORT1=${A[1]}
    NODE_PORT2=${A[2]}
    NODE_PORT3=${A[3]}
    NODE_PORT4=${A[4]}
    PROXY_NODE_PORT0=${A[5]}
    PROXY_NODE_PORT1=${A[6]}
    PROXY_NODE_PORT2=${A[7]}
    PROXY_NODE_PORT3=${A[8]}
    PROXY_NODE_PORT4=${A[9]}
    TESTER_PORT=${A[10]}
else
    NODE_PORT0=$2
    NODE_PORT1=$3
    NODE_PORT2=$4
    NODE_PORT3=$5
    NODE_PORT4=$6
    PROXY_NODE_PORT0=$7
    PROXY_NODE_PORT1=$8
    PROXY_NODE_PORT2=$9
    PROXY_NODE_PORT3=${10}
    PROXY_NODE_PORT4=${11}
    TESTER_PORT=${12}
fi

RAFT_TEST=$GOPATH/bin/rafttest
RAFT_NODE=$GOPATH/bin/raftrunner
PROXY_NODE=$GOPATH/bin/raftproxyrunner
ALL_PORTS="${NODE_PORT0},${NODE_PORT1},${NODE_PORT2},${NODE_PORT3},${NODE_PORT4}"
ALL_PROXY_PORTS="${PROXY_NODE_PORT0},${PROXY_NODE_PORT1},${PROXY_NODE_PORT2},${PROXY_NODE_PORT3},${PROXY_NODE_PORT4}"
##################################################

if [ $# -le 1 ]; then
    echo "All real ports:" ${ALL_PORTS}
    echo "All proxy ports:" ${ALL_PROXY_PORTS}
fi

# Start 5 student raft nodes.
${RAFT_NODE} -myport=${NODE_PORT0} -ports=${ALL_PROXY_PORTS} -retries=15 -N=5 -id=0 & # 2> /dev/null &
RAFT_NODE_PID0=$!
sleep 1

${RAFT_NODE} -myport=${NODE_PORT1} -ports=${ALL_PROXY_PORTS} -retries=15 -N=5 -id=1 & # 2> /dev/null &
RAFT_NODE_PID1=$!
sleep 1

${RAFT_NODE} -myport=${NODE_PORT2} -ports=${ALL_PROXY_PORTS} -retries=15 -N=5 -id=2 & # 2> /dev/null &
RAFT_NODE_PID2=$!
sleep 1

${RAFT_NODE} -myport=${NODE_PORT3} -ports=${ALL_PROXY_PORTS} -retries=15 -N=5 -id=3 & # 2> /dev/null &
RAFT_NODE_PID3=$!
sleep 1

${RAFT_NODE} -myport=${NODE_PORT4} -ports=${ALL_PROXY_PORTS} -retries=15 -N=5 -id=4 & # 2> /dev/null &
RAFT_NODE_PID4=$!
sleep 1

# start all the proxy upon each students' nodes
${PROXY_NODE} -raftport=${NODE_PORT0} -proxyport=${PROXY_NODE_PORT0} -id=0 &
PROXY_NODE_PID0=$!
sleep 1

${PROXY_NODE} -raftport=${NODE_PORT1} -proxyport=${PROXY_NODE_PORT1} -id=1 &
PROXY_NODE_PID1=$!
sleep 1

${PROXY_NODE} -raftport=${NODE_PORT2} -proxyport=${PROXY_NODE_PORT2} -id=2 &
PROXY_NODE_PID2=$!
sleep 1

${PROXY_NODE} -raftport=${NODE_PORT3} -proxyport=${PROXY_NODE_PORT3} -id=3 &
PROXY_NODE_PID3=$!
sleep 1

${PROXY_NODE} -raftport=${NODE_PORT4} -proxyport=${PROXY_NODE_PORT4} -id=4 &
PROXY_NODE_PID4=$!
sleep 1

## Start rafttest.
${RAFT_TEST} -proxyports=${ALL_PROXY_PORTS} -N=5 -t $1

# Kill raft node.
kill -9 ${RAFT_NODE_PID0} 2> /dev/null
kill -9 ${RAFT_NODE_PID1} 2> /dev/null
kill -9 ${RAFT_NODE_PID2} 2> /dev/null
kill -9 ${RAFT_NODE_PID3} 2> /dev/null
kill -9 ${RAFT_NODE_PID4} 2> /dev/null

# Kill proxy
kill -9 ${PROXY_NODE_PID0} 2> /dev/null
kill -9 ${PROXY_NODE_PID1} 2> /dev/null
kill -9 ${PROXY_NODE_PID2} 2> /dev/null
kill -9 ${PROXY_NODE_PID3} 2> /dev/null
kill -9 ${PROXY_NODE_PID4} 2> /dev/null


wait ${RAFT_NODE_PID0} 2> /dev/null
wait ${RAFT_NODE_PID1} 2> /dev/null
wait ${RAFT_NODE_PID2} 2> /dev/null
wait ${RAFT_NODE_PID3} 2> /dev/null
wait ${RAFT_NODE_PID4} 2> /dev/null


wait ${PROXY_NODE_PID0} 2> /dev/null
wait ${PROXY_NODE_PID1} 2> /dev/null
wait ${PROXY_NODE_PID2} 2> /dev/null
wait ${PROXY_NODE_PID3} 2> /dev/null
wait ${PROXY_NODE_PID4} 2> /dev/null