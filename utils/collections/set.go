package collections

import (
	"fmt"
	"strings"
)

type (
	// SetItem represent a item present in a Set.
	SetItem interface{}
	// SetSlice is a slice of items present in a Set.
	SetSlice []SetItem
	// Set is just a set.
	Set map[SetItem]membership

	membership struct{}
)

// Insert adds a value to the set. If the set did not have this value present,
// true is returned. If the set did have this value present, false is returned.
func (s Set) Insert(si SetItem) (ok bool) {
	_, ok = s[si]
	s[si] = membership{}

	return !ok
}

// Contains returns true if the set contains a value.
func (s Set) Contains(si SetItem) (ok bool) {
	_, ok = s[si]

	return
}

// Remove a value from the set. Returns true if the value was present in the
// set.
func (s Set) Remove(si SetItem) (ok bool) {
	_, ok = s[si]
	if ok {
		delete(s, si)
	}

	return
}

// Clear clears the set, removing all values.
func (s *Set) Clear() {
	*s = MakeSet()
}

// Collect returns a slice of all the elements present in the set in arbitrary
// order.
func (s Set) Collect() (ss SetSlice) {
	ss = make([]SetItem, 0, s.Size())

	for k := range s {
		ss = append(ss, k)
	}

	return
}

// Union returns the set of values that are in s or in other, without duplicates.
func (s Set) Union(other Set) (set Set) {
	set = MakeSet()

	for k := range s {
		set[k] = membership{}
	}

	for k := range other {
		set[k] = membership{}
	}

	return
}

// Difference returns the set of values that are in s but not in other.
func (s Set) Difference(other Set) (set Set) {
	set = MakeSet()

	for k := range s {
		if ok := other.Contains(k); !ok {
			set[k] = membership{}
		}
	}

	return
}

// ToString returns a string representation of the set in arbitrary order.
func (s Set) ToString() string {
	r := make([]string, 0, s.Size())
	for k := range s {
		r = append(r, fmt.Sprintf("%v", k))
	}

	return fmt.Sprintf("{%v}", strings.Join(r, ","))
}

// IsEqual returns true if s and other are the same.
func (s Set) IsEqual(other Set) bool {
	return s.SymmetricDifference(other).Size() == 0
}

// SymmetricDifference returns the set of values that are in s or in other but not
// in both.
func (s Set) SymmetricDifference(other Set) (set Set) {
	set = MakeSet()

	for k := range s {
		if ok := other.Contains(k); !ok {
			set[k] = membership{}
		}
	}

	for k := range other {
		if ok := s.Contains(k); !ok {
			set[k] = membership{}
		}
	}

	return
}

// Intersection returns the set of values that are both in s and other.
func (s Set) Intersection(other Set) (set Set) {
	set = MakeSet()

	for k := range s {
		if ok := other.Contains(k); ok {
			set[k] = membership{}
		}
	}

	return
}

// IsEmpty returns true if the set contains no elements.
func (s Set) IsEmpty() bool {
	return len(s) == 0
}

// Size returns the cardinality of a set.
func (s Set) Size() int {
	return len(s)
}

// SubsetOf returns true if s is a subset of other.
func (s Set) SubsetOf(other Set) bool {
	if s.Size() > other.Size() {
		return false
	}

	for k := range s {
		if _, exists := other[k]; !exists {
			return false
		}
	}

	return true
}

// MakeSet creates a new Set.
func MakeSet(si ...SetItem) (s Set) {
	s = make(map[SetItem]membership)
	for _, v := range si {
		s[v] = membership{}
	}

	return
}
