package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cmu440/bitcoin"
	"github.com/cmu440/lsp"
)

// Attempt to connect miner as a client to the server.
func joinWithServer(hostport string) (lsp.Client, error) {
	// TODO: implement this!
	miner, err := lsp.NewClient(hostport, lsp.NewParams())
	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(bitcoin.NewJoin())
	if err := miner.Write(bytes); err != nil {
		return nil, err
	}
	return miner, nil
}

func main() {
	const numArgs = 2
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <hostport>", os.Args[0])
		return
	}

	hostport := os.Args[1]
	miner, err := joinWithServer(hostport)
	if err != nil {
		fmt.Println("Failed to join with server:", err)
		return
	}

	defer miner.Close()

	// TODO: implement this!
	for {
		bytes, err := miner.Read()
		if err != nil {
			fmt.Println("Failed to read from server:", err)
			return
		}

		task := &bitcoin.Message{}
		if err := json.Unmarshal(bytes, task); err != nil {
			fmt.Println("Unable to unmarshal response from server:", err)
			return
		}

		var minHash uint64
		minHash = 18446744073709551615 // set as largest uint64 value first
		nonce := task.Lower
		for n := task.Lower; n <= task.Upper; n++ {
			hash := bitcoin.Hash(task.Data, n)
			if hash < minHash {
				minHash = hash
				nonce = n
			}
		}
		//send minhash and corresponding nonce back to server
		response, _ := json.Marshal(bitcoin.NewResult(minHash, nonce))
		if err := miner.Write(response); err != nil {
			fmt.Println("Failed to send result to server:", err)
			return
		}
	}
}
