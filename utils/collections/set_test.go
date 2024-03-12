package collections_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	coll "github.com/xgodev/boost/utils/collections"
)

// ====================
// Tests
// ====================

func TestSet_ToString(t *testing.T) {
	set := coll.MakeSet("cat", "dog", "cow")

	assert.Contains(t, set.ToString(), "cat", "they should be equal")
	assert.Contains(t, set.ToString(), "dog", "they should be equal")
	assert.Contains(t, set.ToString(), "cow", "they should be equal")
}

func TestSet_IsEqual(t *testing.T) {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "dog")
	set3 := coll.MakeSet("frog")
	set4 := coll.MakeSet("cow", "cat", "dog")

	assert.Equal(t, true, set4.IsEqual(set4), "they should be equal")
	assert.Equal(t, true, set1.IsEqual(set4), "they should be equal")
	assert.Equal(t, false, set2.IsEqual(set3), "they should be equal")
}

func TestSet_Insert(t *testing.T) {
	set := coll.MakeSet()

	ok := set.Insert("cat")
	assert.Equal(t, true, ok, "they should be equal")

	ok = set.Insert("cat")
	assert.Equal(t, false, ok, "they should be equal")
}

func TestSet_SubsetOf(t *testing.T) {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "dog")
	set3 := coll.MakeSet("frog")
	set4 := coll.MakeSet("cat", "dog", "frog", "cow")

	assert.Equal(t, true, set2.SubsetOf(set1), "they should be equal")
	assert.Equal(t, false, set1.SubsetOf(set2), "they should be equal")
	assert.Equal(t, false, set4.SubsetOf(set3), "they should be equal")
	assert.Equal(t, true, set3.SubsetOf(set4), "they should be equal")
}

func TestSet_Contains(t *testing.T) {
	set := coll.MakeSet("cat", "dog", "cow")

	ok := set.Contains("cat")
	assert.Equal(t, true, ok, "they should be equal")

	ok = set.Contains("buffalo")
	assert.Equal(t, false, ok, "they should be equal")
}

func TestSet_Remove(t *testing.T) {
	set := coll.MakeSet("cat", "dog", "cow")

	ok := set.Remove("cat")
	assert.Equal(t, true, ok, "they should be equal")

	ok = set.Remove("cat")
	assert.Equal(t, false, ok, "they should be equal")
}

func TestSet_Intersection(t *testing.T) {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "duck", "bull")

	intersection := set1.Intersection(set2).Collect()
	assert.ElementsMatch(t, intersection, coll.SetSlice{"cat"})
}

func TestSet_SymmetricDifference(t *testing.T) {
	set1 := coll.MakeSet(1, 2, 3)
	set2 := coll.MakeSet(4, 2, 3, 4)

	symDiff := set1.SymmetricDifference(set2).Collect()
	assert.ElementsMatch(t, symDiff, coll.SetSlice{1, 4})
}

func TestSet_Difference(t *testing.T) {
	set1 := coll.MakeSet(1, 2, 3)
	set2 := coll.MakeSet(4, 2, 3, 4)

	diff1 := set1.Difference(set2).Collect()
	assert.ElementsMatch(t, diff1, coll.SetSlice{1})

	diff2 := set2.Difference(set1).Collect()
	assert.ElementsMatch(t, diff2, coll.SetSlice{4})
}

func TestSet_Union(t *testing.T) {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "duck", "bull")

	union := set1.Union(set2).Collect()
	assert.ElementsMatch(t, union, coll.SetSlice{"dog", "cow", "duck", "bull", "cat"})
}

func TestSet_Collect(t *testing.T) {
	els := coll.SetSlice{"cat", "cow", 10, true, false, 10, true, false, "cat", "cow"}

	setValues := coll.MakeSet(els...).Collect()
	assert.ElementsMatch(t, setValues, coll.SetSlice{"cat", 10, "cow", true, false})
}

func TestSet_Clear(t *testing.T) {
	set := coll.MakeSet("cat", "dog", "cow")

	set.Clear()
	assert.Equal(t, 0, set.Size(), "they should be equal")
}

func TestSet_IsEmpty(t *testing.T) {
	set := coll.MakeSet()
	assert.Equal(t, true, set.IsEmpty(), "they should be equal")

	set.Insert("cat")
	assert.Equal(t, false, set.IsEmpty(), "they should be equal")
}

// ====================
// Examples
// ====================

func ExampleSet_Insert() {
	set := coll.MakeSet()

	ok := set.Insert("cat")
	fmt.Println(ok)

	ok = set.Insert("cat")
	fmt.Println(ok)
	// Output:
	// true
	// false
}

func ExampleSet_Contains() {
	set := coll.MakeSet("cat", "dog", "cow")

	ok := set.Contains("cat")
	fmt.Println(ok)

	ok = set.Contains("buffalo")
	fmt.Println(ok)
	// Output:
	// true
	// false
}

func ExampleSet_Remove() {
	set := coll.MakeSet("cat", "dog", "cow")

	ok := set.Remove("cat")
	fmt.Println(ok)

	ok = set.Remove("cat")
	fmt.Println(ok)
	// Output:
	// true
	// false
}

func ExampleSet_Intersection() {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "duck", "bull")
	intersection := set1.Intersection(set2)

	fmt.Println(intersection.Contains("cat"))
	fmt.Println(intersection.Contains("dog"))
	fmt.Println(intersection.Contains("cow"))
	fmt.Println(intersection.Contains("duck"))
	fmt.Println(intersection.Contains("bull"))
	// Output:
	// true
	// false
	// false
	// false
	// false
}

func ExampleSet_SymmetricDifference() {
	set1 := coll.MakeSet(1, 2, 3)
	set2 := coll.MakeSet(4, 2, 3, 4)
	symDiff := set1.SymmetricDifference(set2)

	fmt.Println(symDiff.Contains(1))
	fmt.Println(symDiff.Contains(4))
	fmt.Println(symDiff.Contains(2))
	fmt.Println(symDiff.Contains(3))
	// Output:
	// true
	// true
	// false
	// false
}

func ExampleSet_Difference() {
	set1 := coll.MakeSet(1, 2, 3)
	set2 := coll.MakeSet(4, 2, 3, 4)

	diff1 := set1.Difference(set2)
	fmt.Println(diff1.Contains(1))
	fmt.Println(diff1.Contains(2))
	fmt.Println(diff1.Contains(3))
	fmt.Println(diff1.Contains(4))

	diff2 := set2.Difference(set1)
	fmt.Println(diff2.Contains(4))
	fmt.Println(diff2.Contains(1))
	fmt.Println(diff2.Contains(2))
	fmt.Println(diff2.Contains(3))
	// Output:
	// true
	// false
	// false
	// false
	// true
	// false
	// false
	// false
}

func ExampleSet_Union() {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "duck", "bull")
	union := set1.Union(set2)

	fmt.Println(union.Contains("cat"))
	fmt.Println(union.Contains("dog"))
	fmt.Println(union.Contains("cow"))
	fmt.Println(union.Contains("duck"))
	fmt.Println(union.Contains("bull"))
	// Output:
	// true
	// true
	// true
	// true
	// true
}

func ExampleSet_SubsetOf() {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "dog")
	set3 := coll.MakeSet("frog")
	set4 := coll.MakeSet("cat", "dog", "frog", "cow")

	fmt.Println(set2.SubsetOf(set1))
	fmt.Println(set1.SubsetOf(set2))
	fmt.Println(set4.SubsetOf(set3))
	fmt.Println(set3.SubsetOf(set4))
	// Output:
	// true
	// false
	// false
	// true
}

func ExampleSet_Clear() {
	set := coll.MakeSet("cat", "dog", "cow")

	set.Clear()
	fmt.Println(set.Size())
	// Output:
	// 0
}

func ExampleSet_IsEmpty() {
	set := coll.MakeSet()
	fmt.Println(set.IsEmpty())

	set.Insert("cat")
	fmt.Println(set.IsEmpty())
	// Output:
	// true
	// false
}

func ExampleSet_IsEqual() {
	set1 := coll.MakeSet("cat", "dog", "cow")
	set2 := coll.MakeSet("cat", "dog")
	set3 := coll.MakeSet("frog")
	set4 := coll.MakeSet("cow", "cat", "dog")

	fmt.Println(set4.IsEqual(set4))
	fmt.Println(set1.IsEqual(set4))
	fmt.Println(set2.IsEqual(set3))
	// Output:
	// true
	// true
	// false
}
