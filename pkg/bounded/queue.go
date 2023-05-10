package bounded

type Queue[T any] struct {
	buffer []T
	limit  int // maximum number of entries permitted - should be the same as len(buffer)
	head   int // index of the head element (next to popped)
	count  int // number of elements currently in the queue
}

func NewQueue[T any](limit int) *Queue[T] {
	return &Queue[T]{
		buffer: make([]T, limit),
		limit:  limit,
		count:  0,
	}
}

func (q *Queue[T]) PushBack(item T) (discarded bool) {
	// discard an element from the head of the queue if necessary
	if q.Full() {
		_, _ = q.Pop()
		discarded = true
	}

	idx := (q.head + q.count) % q.limit
	q.buffer[idx] = item
	q.count++
	return
}

func (q *Queue[T]) Pop() (item T, ok bool) {
	if q.Empty() {
		ok = false
		return
	}

	item = q.buffer[q.head]
	q.head = (q.head + 1) % q.limit
	q.count--
	ok = true
	return
}

func (q *Queue[T]) Peek() (item T, ok bool) {
	if q.Empty() {
		ok = false
		return
	}

	item = q.buffer[q.head]
	ok = true
	return
}

func (q *Queue[T]) Empty() bool {
	return q.count == 0
}

func (q *Queue[T]) Full() bool {
	return q.count == q.limit
}
