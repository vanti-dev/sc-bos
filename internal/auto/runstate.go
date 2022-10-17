package auto

// RunState describes the states an automation can be in
//
//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=RunState
type RunState int

const (
	RunStateIdle RunState = iota
	RunStateStarting
	RunStateRunning
	RunStateTransientFailure
	RunStateStopped
	RunStateFailed
)

func (r RunState) IsTerminal() bool {
	switch r {
	case RunStateFailed, RunStateStopped:
		return true
	default:
		return false
	}
}
