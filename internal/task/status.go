package task

type Status string

const (
	StatusInactive Status = "inactive" // Stopped task. Status after calling Stop
	StatusLoading  Status = "loading"  // The task is loading configuration. Status while Configure is running
	StatusActive   Status = "active"   // The task has valid config and is serving requests. Status after calling Start
	StatusError    Status = "error"    // The task failed to start. If start fails, the driver will be in this state.
)
