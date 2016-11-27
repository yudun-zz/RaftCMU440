#!/usr/bin/env bash

rm ./rafttest.log  2> /dev/null

# generate 11 distinct random numbers
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
ALL_PORTS=" ${NODE_PORT0} ${NODE_PORT1} ${NODE_PORT2} ${NODE_PORT3} ${NODE_PORT4}"
ALL_PROXY_PORTS=" ${PROXY_NODE_PORT0} ${PROXY_NODE_PORT1} ${PROXY_NODE_PORT2} ${PROXY_NODE_PORT3} ${PROXY_NODE_PORT4}"

echo "All real ports:" ${ALL_PORTS}
echo "All proxy ports:" ${ALL_PROXY_PORTS}

./rafttest_single.sh testOneCandidateOneRoundElection ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testOneCandidateStartTwoElection ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testTwoCandidateForElection ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testSplitVote1 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testSplitVote2 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testAllForElection1 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testAllForElection2 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testLeaderRevertToFollower1 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testLeaderRevertToFollower2 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testLeaderRevertToFollower3 ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}

./rafttest_single.sh testOneSimplePut ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testOneSimpleUpdate ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testOneSimpleDelete ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}
./rafttest_single.sh testDeleteNonExistKey ${ALL_PORTS}  ${ALL_PROXY_PORTS} ${TESTER_PORT}

cat ./rafttest.log
rm ./rafttest.log