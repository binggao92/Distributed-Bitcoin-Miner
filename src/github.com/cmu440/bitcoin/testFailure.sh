#!/bin/sh
# Client Sends request to server: [Request 6073511289018497495 0 9999999].
# Expecting result:
$GOPATH/bin/server 6060 &
for i in {1..2}
do
    $GOPATH/bin/miner localhost:6060 &
done
ts=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
$GOPATH/bin/client localhost:6060 6073511289018497495 9999999
sleep 0.2
#kill miner process
ps -ef | grep $GOPATH/bin/miner | grep -v grep | awk '{print $2}' | xargs kill
#connect a new miner
sleep 0.2
until $GOPATH/bin/miner localhost:6060 | grep -m 1 ^Result; do : ; done
tt=$(ruby -e 'puts (Time.now.to_f * 1000).to_i')
elapsed=$((tt - ts))
echo "Time taken: $elapsed milliseconds" > testFailure_time.txtsh
ps -ef | grep $GOPATH/bin | grep -v grep | awk '{print $2}' | xargs kill
