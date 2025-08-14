package emergencylightpb

import (
	"context"
	"math/rand"
	"time"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-golang/pkg/resource"
	"github.com/vanti-dev/sc-bos/pkg/gen"
)

type Model struct {
	testResultSet *resource.Value // of *gen.TestResultSet
}

func NewModel(opts ...resource.Option) *Model {
	defaultOpts := []resource.Option{resource.WithInitialValue(&gen.TestResultSet{})}
	opts = append(defaultOpts, opts...)
	return &Model{
		testResultSet: resource.NewValue(opts...),
	}
}

func (m *Model) SetLastDurationTest(result gen.EmergencyTestResult_Result) {
	_, _ = m.testResultSet.Set(&gen.EmergencyTestResult{
		EndTime: timestamppb.Now(),
		Result:  result,
	}, resource.WithUpdateMask(&fieldmaskpb.FieldMask{
		Paths: []string{"duration_test"},
	}))
}

func (m *Model) SetLastFunctionalTest(result gen.EmergencyTestResult_Result) {
	_, _ = m.testResultSet.Set(&gen.EmergencyTestResult{
		EndTime: timestamppb.Now(),
		Result:  result,
	}, resource.WithUpdateMask(&fieldmaskpb.FieldMask{
		Paths: []string{"function_test"},
	}))
}

func (m *Model) GetTestResultSet() *gen.TestResultSet {
	return m.testResultSet.Get().(*gen.TestResultSet)
}

func (m *Model) RunDurationTest() {
	result := &gen.EmergencyTestResult{
		EndTime: timestamppb.Now(),
		Result:  getRandomEmergencyLightResult(),
	}
	m.SetLastDurationTest(result.Result)
}

func (m *Model) RunFunctionTest() {
	result := &gen.EmergencyTestResult{
		EndTime: timestamppb.Now(),
		Result:  getRandomEmergencyLightResult(),
	}
	m.SetLastFunctionalTest(result.Result)
}

func getRandomEmergencyLightResult() gen.EmergencyTestResult_Result {
	n := rand.Intn(11)
	return gen.EmergencyTestResult_Result(n)
}

func (m *Model) PullTestResults(ctx context.Context, opts ...resource.ReadOption) <-chan PullTestResultSetChange {
	send := make(chan PullTestResultSetChange)

	recv := m.testResultSet.Pull(ctx, opts...)
	go func() {
		defer close(send)
		for change := range recv {
			value := change.Value.(*gen.TestResultSet)
			send <- PullTestResultSetChange{
				Value:      value,
				ChangeTime: change.ChangeTime,
			}
		}
	}()

	return send
}

type PullTestResultSetChange struct {
	Value      *gen.TestResultSet
	ChangeTime time.Time
}
