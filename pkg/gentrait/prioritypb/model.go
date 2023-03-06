package prioritypb

import (
	"github.com/vanti-dev/sc-bos/pkg/priority"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Model interface {
	ClearIndex(i int32) error
	ClearName(name string) error
}

func newModel[T any](names ...string) *model[T] {
	return &model[T]{
		list:  priority.NewList[T](len(names)),
		names: names,
	}
}

type model[T any] struct {
	list  *priority.List[T]
	names []string
}

func (m *model[T]) ClearIndex(i int32) error {
	m.list.Clear(int(i))
	return nil
}

func (m *model[T]) ClearName(name string) error {
	i, ok := m.indexForName(name)
	if !ok {
		return status.Error(codes.InvalidArgument, "unknown name")
	}
	m.list.Clear(i)
	return nil
}

func (m *model[T]) indexForName(name string) (int, bool) {
	for i, s := range m.names {
		if s == name {
			return i, true
		}
	}
	return 0, false
}
