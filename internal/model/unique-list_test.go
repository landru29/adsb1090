package model_test

import (
	"sort"
	"testing"

	"github.com/landru29/adsb1090/internal/model"
	"github.com/stretchr/testify/assert"
)

type element string

// Empty implements the model.Uniquer interface.
func (e element) Empty() bool {
	return e.String() == ""
}

// Canonical implements the model.Uniquer interface.
func (e element) Canonical() string {
	return string(e)
}

// String implements the Stringer interface.
func (e element) String() string {
	return string(e)
}

func TestUniqueList(t *testing.T) {
	t.Parallel()

	t.Run("add elements", func(t *testing.T) {
		t.Parallel()

		list := model.UniqueList[element]{}

		list.Add(element("bar"))
		list.Add(element("foo"))
		list.Add(element("foo"))
		list.Add(element(""))

		assert.Equal(t, 2, list.Len())
	})

	t.Run("sort elements", func(t *testing.T) {
		t.Parallel()

		list := model.UniqueList[element]{}

		list.Add(element("bar"))
		list.Add(element("foo"))
		list.Add(element("foo"))

		sort.Sort(list)

		assert.Equal(t, "foo", list.First().String())
	})
}
