package model

import (
	"strings"
)

// Uniquer is a unique element of a list.
type Uniquer interface {
	comparable
	Empty() bool
	Canonical() string
}

// UniqueList is a list that ensure that all element are unique.
type UniqueList[T Uniquer] struct {
	data []T
}

// Add adds non-empty elements in the list.
func (u *UniqueList[T]) Add(str T) {
	if str.Empty() {
		return
	}

	for _, elt := range u.data {
		if elt == str {
			return
		}
	}

	u.data = append(u.data, str)
}

// Len imlplements the sort.Interface interface.
func (u UniqueList[T]) Len() int {
	return len(u.data)
}

// Swap imlplements the sort.Interface interface.
func (u UniqueList[T]) Swap(i, j int) {
	u.data[i], u.data[j] = u.data[j], u.data[i]
}

// Less imlplements the sort.Interface interface.
func (u UniqueList[T]) Less(i, j int) bool {
	return strings.Compare(u.data[i].Canonical(), u.data[j].Canonical()) > 0
}

// First is the first element of the list.
func (u UniqueList[T]) First() T { //nolint: ireturn
	return u.data[0]
}
