package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cmu440/bitcoin"
	"github.com/cmu440/lsp"
)

type server struct {
	lspServer  lsp.Server
	jobQueue   *list.List       // hold jobs from client
	miners     map[int]*job     // record miners information
	clients    map[int]*request // store the information about requests from clients
	receiveMsg chan idMsg       // a channel to receive msg
	connLost   chan int         // a channel to check connection lost
}

type job struct {
	clientid int
	message  string
	minNonce uint64
	maxNonce uint64
}

//store the details about the request from client
type request struct {
	jobsRemain int
	minHash    uint64
	nonce      uint64
}

type idMsg struct {
	id  int
	msg *bitcoin.Message
}

const maxJobSize = 2000

var LOGF *log.Logger

func startServer(port int) (*server, error) {
	// TODO: implement this!
	lspserver, err := lsp.NewServer(port, lsp.NewParams())
	if err != nil {
		return nil, err
	}
	srv := &server{}
	srv.lspServer = lspserver
	srv.jobQueue = list.New()
	srv.miners = make(map[int]*job)
	srv.clients = make(map[int]*request)
	srv.receiveMsg = make(chan idMsg)
	srv.connLost = make(chan int)
	return srv, nil
}

func main() {
	// You may need a logger for debug purpose
	const (
		name = "log.txt"
		flag = os.O_RDWR | os.O_CREATE
		perm = os.FileMode(0666)
	)

	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return
	}
	defer file.Close()

	LOGF = log.New(file, "", log.Lshortfile|log.Lmicroseconds)
	// Usage: LOGF.Println() or LOGF.Printf()
	LOGF.Println("Test Starts here")

	const numArgs = 2
	if len(os.Args) != numArgs {
		fmt.Printf("Usage: ./%s <port>", os.Args[0])
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Port must be a number:", err)
		return
	}

	srv, err := startServer(port)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Server listening on port", port)

	defer srv.lspServer.Close()

	// TODO: implement this!
	go srv.receiveMessage()

	srv.handleMessage()
	file.Close()
}

func (s *server) receiveMessage() {
	for {
		id, bytes, err := s.lspServer.Read()
		if err != nil {
			s.connLost <- id
		} else {
			msg := &bitcoin.Message{}
			err := json.Unmarshal(bytes, msg)
			if err != nil {
				continue
			}
			s.receiveMsg <- idMsg{id, msg}
		}
	}
}

func (s *server) handleMessage() {
	for {
		select {
		case idAndMsg := <-s.receiveMsg: //when receive msg
			id := idAndMsg.id
			msg := idAndMsg.msg
			switch msg.Type {
			case bitcoin.Join:
				s.miners[id] = nil // Add a miner to the miner pool
			case bitcoin.Request:
				s.addJob(id, msg)
			case bitcoin.Result:
				clientId := s.miners[id].clientid
				s.miners[id] = nil
				clientJob, exist := s.clients[clientId]
				if exist {
					clientJob.jobsRemain--
					if msg.Hash < s.clients[clientId].minHash {
						clientJob.minHash = msg.Hash
						clientJob.nonce = msg.Nonce
					}
					// Client's job finished. send result back to client
					if clientJob.jobsRemain == 0 {
						s.sendMessage(clientId, bitcoin.NewResult(clientJob.minHash, clientJob.nonce))
						LOGF.Println("Job finished:", clientId, clientJob.minHash, clientJob.nonce)
						s.lspServer.CloseConn(clientId) //close the link between client and server
						delete(s.clients, clientId)
					}
				}
			}
		case id := <-s.connLost: //when connection lost
			if _, exist := s.clients[id]; exist {
				delete(s.clients, id)
			} else if job, exist := s.miners[id]; exist {
				if job != nil {
					s.jobQueue.PushFront(*job)
				}
				delete(s.miners, id)
			}
		}
		s.dispatchJobsToMiners()
	}

}

func (s *server) sendMessage(id int, msg *bitcoin.Message) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err := s.lspServer.Write(id, bytes); err != nil {
		return err
	}
	return nil
}

func (s *server) addJob(id int, task *bitcoin.Message) {
	newRequest := &request{0, ^uint64(0), uint64(0)}
	lower, upper := task.Lower, task.Lower+maxJobSize-1

	for upper < task.Upper {
		newRequest.jobsRemain++
		s.jobQueue.PushBack(job{id, task.Data, lower, upper})
		lower = upper + 1
		upper = lower + maxJobSize - 1
	}
	newRequest.jobsRemain++
	s.jobQueue.PushBack(job{id, task.Data, lower, task.Upper})
	s.clients[id] = newRequest
}

func (s *server) findIdleMiner() (int, bool) {
	for minerId, job := range s.miners {
		if job == nil {
			LOGF.Println("A idle miner(Id is", minerId, ") starts working")
			return minerId, true
		}
	}
	return 0, false
}

func (s *server) dispatchJobsToMiners() {
	for s.jobQueue.Len() > 0 {
		job := s.jobQueue.Front().Value.(job)
		if _, exist := s.clients[job.clientid]; exist {
			if minerId, exist := s.findIdleMiner(); exist {
				s.miners[minerId] = &job
				s.sendMessage(minerId, bitcoin.NewRequest(job.message, job.minNonce, job.maxNonce))
			} else {
				return
			}
		}
		s.jobQueue.Remove(s.jobQueue.Front())
	}
}
