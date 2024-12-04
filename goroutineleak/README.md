## Goroutine Leak

### 1. Start a dev server

```bash
temporal server start-dev
```

### 2. From the root of the project, start a Worker

```bash
go run goroutineleak/worker/main.go
```

### 3. Start the Workflow Execution

```bash
go run goroutineleak/starter/main.go
```

### 4. Inspect number of goroutines

Wait for the workflow from step 3 to complete before downloading a profile.

```bash
go tool pprof -top http://localhost:6060/debug/pprof/goroutine
```

Output should look something like:

```
File: main
Type: goroutine
Time: Dec 4, 2024 at 5:39pm (EST)
Showing nodes accounting for 992, 99.70% of 995 total
Dropped 64 nodes (cum <= 4)
      flat  flat%   sum%        cum   cum%
       992 99.70% 99.70%        992 99.70%  runtime.gopark
         0     0% 99.70%          4   0.4%  github.com/grpc-ecosystem/go-grpc-middleware/retry.UnaryClientInterceptor.func1
         0     0% 99.70%        966 97.09%  github.com/temporalio/samples-go/goroutineleak.SampleChildWorkflow
```

Repeat steps 3-4 and observe that number of goroutines continues to climb, while number of running [workflows in the UI](http://localhost:8233/namespaces/default/workflows) remains zero.

## Observations
- Goroutines max out at ~10,000 when no worker options are set. This lines up with the default sticky workflow cache size; after setting cache size to 2k goroutines max out around 2,000.
- Child workflows do not have to be started from inside `workflow.Go` in order to leak goroutines.
- Changing the parent close policy from "terminate" to "request cancel" in child workflow options avoids leaking goroutines.
- Workflow functions that simply sleep for 30 seconds still have open goroutines more than 30 seconds after being terminated.

## Open Questions
- Is this related to https://github.com/temporalio/sdk-go/issues/716 or https://github.com/temporalio/sdk-go/issues/1049?
