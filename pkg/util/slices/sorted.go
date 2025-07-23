package slices

import (
	"cmp"
	"slices"
	"sort"
)

// Sorted is a slice of unique items that are kept in sorted order.
// Changing the compared values of items in the slice will not automatically reorder the slice.
type Sorted[T any] struct {
	items []T
	cmp   func(a, b T) int
}

// NewSorted is like NewSortedFunc(cmp.Compare[T], items...).
func NewSorted[T cmp.Ordered](items ...T) *Sorted[T] {
	return NewSortedFunc(cmp.Compare[T], items...)
}

// NewSortedFunc creates a new Sorted slice with the given items and comparison function.
// The items slice is sorted in place and not copied, do not modify it after passing it to NewSortedFunc.
func NewSortedFunc[T any](cmp func(a, b T) int, items ...T) *Sorted[T] {
	if len(items) > 0 {
		sort.Slice(items, func(i, j int) bool {
			return cmp(items[i], items[j]) < 0
		})
	}
	return &Sorted[T]{items: items, cmp: cmp}
}

// Cmp returns the underlying comparison function.
func (s *Sorted[T]) Cmp() func(a, b T) int {
	return s.cmp
}

// Set adds or replaces an item to the slice, returning the index, old item, and true if the item was replaced.
func (s *Sorted[T]) Set(item T) (int, T, bool) {
	i, found := s.find(item)
	if found {
		old := s.items[i]
		s.items[i] = item
		return i, old, true
	}
	// insert item at index i
	s.items = append(s.items, item) // grow the slice
	if i < s.Len()-1 {
		copy(s.items[i+1:], s.items[i:s.Len()-1])
		s.items[i] = item
	}
	var zero T
	return i, zero, false
}

// Remove removes an item from the slice, returning the index and true if the item was removed.
func (s *Sorted[T]) Remove(item T) (int, T, bool) {
	i, found := s.find(item)
	var zero T
	if found {
		removed := s.items[i]
		s.items = append(s.items[:i], s.items[i+1:]...)
		return i, removed, true
	}
	return 0, zero, false
}

// Len returns the number of items in the slice.
func (s *Sorted[T]) Len() int {
	return len(s.items)
}

// All calls the given function for each item in the slice, breaking early if the function returns false.
func (s *Sorted[T]) All(yield func(int, T) bool) {
	for i, item := range s.items {
		if !yield(i, item) {
			break
		}
	}
}

// Get returns the item at the given index.
func (s *Sorted[T]) Get(i int) T {
	return s.items[i]
}

// Find returns smallest index in s where item would exist.
// If the existing item at the index compares equal to item, then it is returned along with true.
func (s *Sorted[T]) Find(item T) (int, T, bool) {
	i, found := s.find(item)
	if found {
		return i, s.items[i], true
	}
	var zero T
	return i, zero, false
}

// Sort sorts the slice in place.
// Set and Remove will maintain the order of the slice, but if the items are changed externally in a way that would change the order, you must call Sort.
func (s *Sorted[T]) Sort() {
	sort.Slice(s.items, func(i, j int) bool {
		return s.cmp(s.items[i], s.items[j]) < 0
	})
}

// Clear removes all items from the slice.
func (s *Sorted[T]) Clear() {
	// todo: use builtin clear once we update to go 1.21
	var zero T
	for i := range s.items {
		s.items[i] = zero // allow GC
	}
	s.items = s.items[:0] // maintain the capacity
}

// Copy returns a shallow copy of s.
func (s *Sorted[T]) Copy() *Sorted[T] {
	return &Sorted[T]{
		items: append([]T(nil), s.items...),
		cmp:   s.cmp,
	}
}

func (s *Sorted[T]) find(item T) (int, bool) {
	return slices.BinarySearchFunc(s.items, item, s.cmp)
}
