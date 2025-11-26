package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sanjayJ369/raft-consensus/internal/node"
	simpletimer "github.com/sanjayJ369/raft-consensus/internal/simpleTimer"
	simpletransport "github.com/sanjayJ369/raft-consensus/internal/simpleTransport"
	"github.com/sanjayJ369/raft-consensus/internal/types"
	"github.com/sanjayJ369/raft-consensus/logger"
)

func main() {
	done := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// simple config
	config := node.Config{
		HeartBeatTimeout:   10 * time.Millisecond,
		ElectionTimeoutMin: 5000 * time.Millisecond,
		ElectionTimeoutMax: 10000 * time.Millisecond,
	}

	// let me starting with 3 nodes
	nodes := make(map[types.NodeId]*node.Node)
	network := make(simpletransport.Network)

	lgr1, close1 := logger.NewLoggerFile("./logs/node1.log", true)
	defer close1()

	lgr2, close2 := logger.NewLoggerFile("./logs/node2.log", true)
	defer close2()

	lgr3, close3 := logger.NewLoggerFile("./logs/node3.log", true)
	defer close3()

	loggers := []types.Logger{
		lgr1, lgr2, lgr3,
	}

	// create nodes
	for i := range 3 {
		nodeId := types.NodeId(i)
		timer := simpletimer.NewSimpleTimer()
		// logger := logger.NewLogger(os.Stdout)

		transport := simpletransport.NewSimpleTransport(nodeId, &network, loggers[i])
		node := node.NewNode(nodeId, nil, config, timer, transport, loggers[i])
		network.Register(nodeId, node)

		// add nodes in cluster as peers for new node
		// add new node as peer to all the prev nodes
		for _, n := range nodes {
			node.AddPeer(n.Id)
			n.AddPeer(nodeId)
		}

		nodes[node.Id] = node
	}

	// start nodes
	for _, node := range nodes {
		go node.EnterFollower()
	}

	// flush logs to disk
	go func() {
		<-signalChan
		for _, lgr := range loggers {
			lgr.Sync()
		}
		os.Exit(1)
	}()

	<-done
}
