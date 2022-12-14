package runstate

// RunState describes the states an automation can be in
//
//go:generate go run golang.org/x/tools/cmd/stringer@latest -type=RunState
type RunState int

const (
	Idle RunState = iota
	Starting
	Running
	TransientFailure
	Stopped
	Failed
)

func (r RunState) IsTerminal() bool {
	switch r {
	case Failed, Stopped:
		return true
	default:
		return false
	}
}
