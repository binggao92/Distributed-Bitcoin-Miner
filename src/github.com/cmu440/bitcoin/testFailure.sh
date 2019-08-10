#!/bin/sh
# Client Sends request to server: [Request 6073511289018497495 0 9999999].
# Expecting result:
$GOPATH/bin/server 6060 &
$GOPATH/bin/miner localhost:6060 &
$GOPATH/bin/miner localhost:6060 &
ts=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
$GOPATH/bin/client localhost:6060 6073511289018497495 999999
sleep 0.1
#kill all miners
ps -ef | grep $GOPATH/bin/miner | grep -v grep | awk '{print $2}' | xargs kill
#connect a new miner
sleep 0.2
$GOPATH/bin/miner localhost:6060

tt=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
elapsed=$((tt - ts))
echo "Time taken: $elapsed milliseconds" > testFailure_time.txtsh
ps -ef | grep $GOPATH/bin/miner | grep -v grep | awk '{print $2}' | xargs kill
ps -ef | grep $GOPATH/bin/server | grep -v grep | awk '{print $2}' | xargs kill
