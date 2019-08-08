# Distributed Bitcoin Miner

Uvic CSC462 Project

Code frame provided by CMU440: https://github.com/cmu440/p1

## Part B(Distributed bitcoin miner)

#### The overall operation of the system should proceed as follows:

• The server is started using the following command, specifying the port to listen on:    

	./server port
• One or more miners are started using the following command, specifying the server’s address and port number separated by a colon:    

	./miner host:port
• When a miner starts, it sends a join request message to the server, letting the server know that it is ready to receive and execute new job requests. New miners may start up and send join requests to the server at any time.
• The user makes a request to the server by executing the following command, speci- fying the server’s address and port number, the message to hash, and the maximum nonce to check:    

	./client host:port message maxNonce
The client program should generate and send a request message to the server, specifying lower nonce “0” and maximum nonce “maxNonce”:    

	[Request message 0 maxNonce]
• The server breaks this request into more manageable-sized jobs and farms them out to its available miners (it is up to you to choose a suitable maximum job size). Once the miner exhausts all possible nonces, it determines the least hash value and its corresponding nonce and sends back the final result:    

	[Result minHash nonce]
• The server collects all final results from the workers, determines the least hash value and its corresponding nonce, and sends back the final result to the request client.

The request client should print the result of each request to standard output as follows (note that you must match this format precisely in order for our automated tests to work):
• If the server responds with a final result, it should print,    

	Result minHash nonce
where minHash and nonce are the values of the lowest hash and its corresponding nonce, respectively.

• If the request client loses its connection with the server, it should simply print Disconnected    

	Result minHash nonce


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
There are serveral test scenarios: test1.sh test2.sh .....
run a test scenario using this command    

    sh ./test1.sh >test1_log.txt
And there will be 2 log files generated: test1_log.txt(all the stdout) and test1_time.txt(total running time)
