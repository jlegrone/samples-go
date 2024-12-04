package goroutineleak

import (
	"time"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/workflow"
)

func SampleParentWorkflow(ctx workflow.Context) (string, error) {
	var futures []workflow.ChildWorkflowFuture
	for i := 0; i < 1000; i++ {
		fut := workflow.ExecuteChildWorkflow(workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{
			WorkflowExecutionTimeout: time.Hour,
			//ParentClosePolicy: enumspb.PARENT_CLOSE_POLICY_REQUEST_CANCEL, // does not leak goroutines
			ParentClosePolicy: enumspb.PARENT_CLOSE_POLICY_TERMINATE, // default; does leak goroutines
		}), SampleChildWorkflow, "World")
		futures = append(futures, fut)
	}

	// Ensure that child workflows have been started
	for _, fut := range futures {
		if err := fut.GetChildWorkflowExecution().Get(ctx, nil); err != nil {
			return "", err
		}
	}

	return "done!", nil
}
