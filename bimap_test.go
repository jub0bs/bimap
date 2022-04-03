package bimap

import (
	"math"
	"sort"
	"testing"

	"golang.org/x/exp/slices"
)

func TestThatANewBimapHasSizeZero(t *testing.T) {
	bi := New[int, string]()
	want := 0
	got := bi.Size()
	if got != want {
		t.Errorf("got %d; want %d", got, want)
	}
}

func TestThatABimapRetainsASingleAssociation(t *testing.T) {
	bi := New[int, string]()
	key := 1
	value := "one"
	bi.Store(key, value)
	v, exists := bi.LoadValue(key)
	if !exists && v != value {
		t.Errorf("got %q; want %q", v, value)
	}
	k, exists := bi.LoadKey(value)
	if !exists && k != key {
		t.Errorf("got %d; want %d", k, key)
	}
}

func TestThatTheZeroValueIsReadilyUsable(t *testing.T) {
	bi := new(Bimap[int, string])
	bi.Store(1, "one")
	v, exists := bi.LoadValue(1)
	if !exists || v != "one" {
		t.Errorf("got %q, %t; want %q, %t", v, exists, "one", true)
	}
}

func TestThatStoringSomeNewPairThatConflictsWithPreviousOnesDiscardsThem(t *testing.T) {
	bi := New[int, string]()
	bi.Store(1, "one")
	bi.Store(2, "two")
	bi.Store(1, "two")
	if size := bi.Size(); size != 1 {
		t.Errorf("bi.Size() = %d; want %d", size, 1)
	}
	if v, exists := bi.LoadValue(1); !exists || v != "two" {
		t.Errorf("got %q, %t; want %q, %t", v, exists, "two", true)
	}
}

func TestDeleteByKeyRemovesTheCorrespondingKeyValuePair(t *testing.T) {
	bi := New[int, string]()
	key := 1
	value := "one"
	bi.Store(key, value)
	bi.DeleteByKey(key)
	if size := bi.Size(); size != 0 {
		t.Errorf("bi.Size() = %d; want %d", size, 0)
	}
	if v, exists := bi.LoadValue(1); exists {
		t.Errorf("got %q, %t; want %q, %t", v, exists, "", false)
	}
}

func TestDeleteByValueRemovesTheCorrespondingKeyValuePair(t *testing.T) {
	bi := New[int, string]()
	key := 1
	value := "one"
	bi.Store(key, value)
	bi.DeleteByValue(value)
	if size := bi.Size(); size != 0 {
		t.Errorf("bi.Size() = %d; want %d", size, 0)
	}
	if k, exists := bi.LoadKey("two"); exists {
		t.Errorf("got %d, %t; want %d, %t", k, exists, 0, false)
	}
}

func TestKeysReturnsAllTheKeysInTheBimap(t *testing.T) {
	bi := New[int, string]()
	bi.Store(1, "one")
	bi.Store(2, "two")
	bi.Store(3, "three")
	got := bi.Keys()
	sort.Ints(got)
	want := []int{1, 2, 3}
	if ok := slices.Equal(got, want); !ok {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestValuesReturnsAllTheValuesInTheBimap(t *testing.T) {
	bi := New[int, string]()
	bi.Store(1, "one")
	bi.Store(2, "two")
	bi.Store(3, "three")
	got := bi.Values()
	sort.Strings(got)
	want := []string{"one", "two", "three"}
	sort.Strings(want)
	if ok := slices.Equal(got, want); !ok {
		t.Errorf("got %v; want %v", got, want)
	}
}

func TestThatNonReflexiveKeyIsRejected(t *testing.T) {
	bi := New[float64, string]()
	ok := bi.Store(math.NaN(), "NaN")
	if size := bi.Size(); ok || size != 0 {
		t.Errorf("got %v, %d; want false, 0", ok, size)
	}
}

func TestThatNonReflexiveValueIsRejected(t *testing.T) {
	bi := New[string, float64]()
	ok := bi.Store("NaN", math.NaN())
	if size := bi.Size(); ok || size != 0 {
		t.Errorf("got %v, %d; want false, 0", ok, size)
	}
}

func TestString(t *testing.T) {
	bi := New[int, string]()
	bi.Store(1, "one")
	bi.Store(2, "two")
	bi.Store(3, "three")
	got := bi.String()
	want := "Bimap[1:one 2:two 3:three]"
	if got != want {
		t.Errorf("got %v; want %v", got, want)
	}
}
