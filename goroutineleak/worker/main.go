package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/goroutineleak"
)

// @@@SNIPSTART samples-go-child-workflow-example-worker-starter
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Validate that goroutines is limited to this number
	worker.SetStickyWorkflowCacheSize(2000)

	// The client is a heavyweight object that should be created only once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "child-workflow", worker.Options{})

	w.RegisterWorkflow(goroutineleak.SampleParentWorkflow)
	w.RegisterWorkflow(goroutineleak.SampleChildWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

// @@@SNIPEND
