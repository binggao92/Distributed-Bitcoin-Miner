#!/bin/sh
SECONDS=0
$GOPATH/bin/server 6060 &
for i in 1 2 3 4 5
do
    $GOPATH/bin/miner localhost:6060 &
done
until $GOPATH/bin/client localhost:6060 test1 9999 | grep -m 1 ^Result; do sleep 0.2 ; done
duration=$SECONDS
echo "$(($duration / 60)) minutes and $(($duration % 60)) seconds elapsed." > test1_time.txt
ps -ef | grep /Users/binggao/Documents/GitHub/Distributed-Bitcoin-Miner/bin | grep -v grep | awk '{print $2}' | xargs kill
