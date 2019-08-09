#!/bin/sh
# Client Sends request to server: [Request 4598353849669926001 0 99999]. There is only 1 miner avaliable
# Expecting result: [Result 446491093039252 39168]
$GOPATH/bin/server 6060 &
$GOPATH/bin/miner localhost:6060 &
ts=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
#client send job to server and wait until result come out
until $GOPATH/bin/client localhost:6060 4598353849669926001 99999 | grep -m 1 ^Result; do sleep 0.2 ; done
tt=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
elapsed=$((tt - ts))
echo "Time taken: $elapsed milliseconds" > test11_time.txt
ps -ef | grep $GOPATH/bin | grep -v grep | awk '{print $2}' | xargs kill
