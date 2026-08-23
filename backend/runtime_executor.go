package backend

import "fmt"

type RuntimeExecutor interface {
	Apply(service LiveService, desired ServiceRuntimeState) error
}

type DryRunExecutor struct{}

func (DryRunExecutor) Apply(service LiveService, desired ServiceRuntimeState) error {
	if service.ID == "" || service.NodeID == "" || !service.Authorized { return fmt.Errorf("unauthorized or incomplete service") }
	if desired != ServiceLive && desired != ServiceDegraded && desired != ServiceStopped { return fmt.Errorf("unsupported desired state") }
	return nil
}

type ExecutionResult struct {
	ServiceID string
	Desired   ServiceRuntimeState
	Applied   bool
}

func ExecuteService(executor RuntimeExecutor, service LiveService, desired ServiceRuntimeState) (ExecutionResult, error) {
	if executor == nil { return ExecutionResult{}, fmt.Errorf("runtime executor is required") }
	if err := executor.Apply(service, desired); err != nil { return ExecutionResult{}, err }
	return ExecutionResult{ServiceID: service.ID, Desired: desired, Applied: true}, nil
}
