package emergencylightpb

import (
	"context"
	"math/rand"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smart-core-os/sc-bos/pkg/gen"
	"github.com/smart-core-os/sc-bos/pkg/util/resources"
	"github.com/smart-core-os/sc-golang/pkg/resource"
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
	_, _ = m.testResultSet.Set(&gen.TestResultSet{
		DurationTest: &gen.EmergencyTestResult{
			EndTime: timestamppb.Now(),
			Result:  result,
		}}, resource.WithUpdateMask(&fieldmaskpb.FieldMask{
		Paths: []string{"duration_test"},
	}))
}

func (m *Model) SetLastFunctionalTest(result gen.EmergencyTestResult_Result) {
	_, _ = m.testResultSet.Set(&gen.TestResultSet{
		FunctionTest: &gen.EmergencyTestResult{
			EndTime: timestamppb.Now(),
			Result:  result,
		}}, resource.WithUpdateMask(&fieldmaskpb.FieldMask{
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
	return resources.PullValue[*gen.TestResultSet](ctx, m.testResultSet.Pull(ctx, opts...))
}

type PullTestResultSetChange = resources.ValueChange[*gen.TestResultSet]
