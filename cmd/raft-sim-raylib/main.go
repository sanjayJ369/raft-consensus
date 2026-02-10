package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"syscall"

	_ "net/http/pprof" // <-- enables /debug/pprof

	rl "github.com/gen2brain/raylib-go/raylib"
	graphics "github.com/sanjayJ369/raft-consensus/internal/impl/graphics"
	"github.com/sanjayJ369/raft-consensus/logger"
)

const (
	WIN_HEIGHT = 800
	WIN_WIDTH  = 1200
)

func main() {

	// ---------------------------------------------------------
	// 1️⃣ Enable Mutex + Block Profiling
	// ---------------------------------------------------------
	runtime.SetMutexProfileFraction(1) // capture all mutex contention
	runtime.SetBlockProfileRate(1)     // capture blocking events

	// ---------------------------------------------------------
	// 2️⃣ Start CPU Profiling
	// ---------------------------------------------------------
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(cpuFile)
	defer func() {
		pprof.StopCPUProfile()
		cpuFile.Close()
	}()

	// ---------------------------------------------------------
	// 3️⃣ Start Execution Trace
	// ---------------------------------------------------------
	traceFile, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	trace.Start(traceFile)
	defer func() {
		trace.Stop()
		traceFile.Close()
	}()

	// ---------------------------------------------------------
	// 4️⃣ Start pprof HTTP Server (live profiling)
	// ---------------------------------------------------------
	go func() {
		log.Println("pprof running at http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// ---------------------------------------------------------
	// 5️⃣ Normal App Setup
	// ---------------------------------------------------------

	done := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	rl.InitWindow(WIN_WIDTH, WIN_HEIGHT, "Raft Consensus Simulation")

	// --- LOGGING SETUP START ---

	// 1. Manager Logger
	mgrLogger, closeMgr := logger.NewLoggerFile("./logs/manager.log", true)
	defer closeMgr()

	// 2. Node 1 Logger
	node1Logger, closeNode1 := logger.NewLoggerFile("./logs/node1.log", true)
	defer closeNode1()

	// 3. Node 2 Logger
	node2Logger, closeNode2 := logger.NewLoggerFile("./logs/node2.log", true)
	defer closeNode2()

	// 4. Node 3 Logger
	node3Logger, closeNode3 := logger.NewLoggerFile("./logs/node3.log", true)
	defer closeNode3()

	// --- LOGGING SETUP END ---

	// Initialize components with their specific loggers
	manager := graphics.NewManager(mgrLogger)

	node1 := graphics.NewNodeGUI(WIN_WIDTH, WIN_HEIGHT, node1Logger)
	node2 := graphics.NewNodeGUI(WIN_WIDTH, WIN_HEIGHT, node2Logger)
	node3 := graphics.NewNodeGUI(WIN_WIDTH, WIN_HEIGHT, node3Logger)

	manager.RegisterNode(node1)
	manager.RegisterNode(node2)
	manager.RegisterNode(node3)

	node1.Node.AddPeer(node2.Node.ID)
	node1.Node.AddPeer(node3.Node.ID)

	node2.Node.AddPeer(node1.Node.ID)
	node2.Node.AddPeer(node3.Node.ID)

	node3.Node.AddPeer(node1.Node.ID)
	node3.Node.AddPeer(node2.Node.ID)

	manager.StartNodes()

	rl.SetTargetFPS(30)

	// Graceful shutdown handler
	go func() {
		<-signalChan
		rl.CloseWindow()
		done <- true
	}()

	for !rl.WindowShouldClose() {
		manager.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		manager.Render()
		rl.EndDrawing()
	}

	select {
	case done <- true:
	default:
	}
}
