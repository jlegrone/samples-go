package goroutineleak

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func SampleChildWorkflow(ctx workflow.Context, name string) (string, error) {
	if err := workflow.Sleep(ctx, 30*time.Second); err != nil {
		return "", err
	}

	return "Hello " + name + "!", nil
}
