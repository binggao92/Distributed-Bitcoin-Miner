#!/bin/sh
# Client Sends request to server: [Request 6073511289018497495 0 999999]. There is only 1 miner avaliable
# Expecting result: [Result 51022575622139 3899]
$GOPATH/bin/server 6060 &
$GOPATH/bin/miner localhost:6060 &
ts=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
#client send job to server and wait until result come out
until $GOPATH/bin/client localhost:6060 6073511289018497495 999999 | grep -m 1 ^Result; do sleep 0.2 ; done
tt=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
elapsed=$((tt - ts))
echo "Time taken: $elapsed milliseconds" > test21_time.txt
ps -ef | grep $GOPATH/bin | grep -v grep | awk '{print $2}' | xargs kill