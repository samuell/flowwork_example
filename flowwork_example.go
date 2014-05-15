package main

import (
	"github.com/samuell/blow"
	"github.com/samuell/flowwork"
	"github.com/trustmaster/goflow"
	"runtime"
)

const (
	NUMTHREADS = 2
)

// ---------------------------------------------------------
// An Example FlowWork network
// ---------------------------------------------------------
type FlowWorkApp struct {
	flow.Graph
}

func NewFlowWorkApp() *FlowWorkApp {
	network := new(FlowWorkApp)
	network.InitGraphState()

	// Add components
	network.Add(new(flowwork.CommandExecutor), "command_executor")
	network.Add(new(blow.Printer), "printer")

	// Connect components
	network.MapInPort("In", "command_executor", "Command")
	network.Connect("command_executor", "CommandOutput", "printer", "Line")

	return network
}

// ---------------------------------------------------------
// Main method
// ---------------------------------------------------------
var finish chan bool

// Use this handler to let main() know when the network terminates
func (a *FlowWorkApp) Finish() {
	finish <- true
}

func main() {
	// Set the number of Operating System-threads to use
	runtime.GOMAXPROCS(NUMTHREADS)

	// Termination signal channel
	finish = make(chan bool)
	// Create network
	net := NewFlowWorkApp()

	// Create the "In" channel
	in := make(chan []byte)
	net.SetInPort("In", in)

	// Run net
	flow.RunNet(net)

	// Give the Filename to read
	in <- []byte("ls -l")

	close(in)
	<-finish
}
