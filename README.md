# Distributed Bitcoin Miner

Uvic CSC462 Project

Code frame provided by CMU440: https://github.com/cmu440/p1

#Part B(Distributed bitcoin miner)

## Compiling the client, miner & server programs
<pre>

To compile the client, miner, and server programs, use the go install command as follows (these instructions assume your GOPATH is pointing to the project's root p1/ directory):

# Compile the client, miner, and server programs. The resulting binaries
# will be located in the $GOPATH/bin directory.
go install github.com/cmu440/bitcoin/client
go install github.com/cmu440/bitcoin/miner
go install github.com/cmu440/bitcoin/server

# Start the server, specifying the port to listen on.
$GOPATH/bin/server 6060

# Start a miner, specifying the server's host:port.
$GOPATH/bin/miner localhost:6060

# Start the client, specifying the server's host:port, the message
# "bradfitz", and max nonce 9999.
$GOPATH/bin/client localhost:6060 bradfitz 9999

</pre>

## Running the tests
there are serveral test scenarios: test1.sh test2.sh .....
run a test scenario using this command
    sh ./test1.sh >test1_log.txt
And there will be 2 log files: test1_log.txt(all the stdout) and test1_time.txt(total running time)
