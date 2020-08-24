package errors

import (
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

// init rand.Seed to make sure tests always produce the same output to avoid any flakiness
func TestMain(m *testing.M) {
	rand.Seed(7894312849)
	os.Exit(m.Run())
}

func Test_KVs_GivesBasicKeyValuePairs(t *testing.T) {
	e := New("hello", "a", 1)
	actual := KVs(e)
	expected := []interface{}{"a", 1}
	if !compareKVs(actual, expected) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func Test_KVs_TraversesKeyValuePairs(t *testing.T) {
	e1 := New("hello", "a", 1)
	e2 := Wrap(e1, "hello", "b", 2)
	actual := KVs(e2)
	expected := []interface{}{"a", 1, "b", 2}
	compareKVs(actual, expected)
}

func Test_KVs_TraversesTheEntireChain(t *testing.T) {
	e1 := randError()
	e2 := Wrap(e1, randString(), randomKeyValuePairs()...)
	e3 := Wrap(e2, randString(), randomKeyValuePairs()...)
	actual := KVs(e3)
	expected := append(e1.keyValuePairs, append(e2.keyValuePairs, e3.keyValuePairs...)...)
	if !compareKVs(actual, expected) {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func Test_KVs_ReturnsNilWhenPassedNil(t *testing.T) {
	kvs := KVs(nil)
	if kvs != nil {
		t.Fatalf("expected nil, got %+v", kvs)
	}
}

func Test_Unwrap_TraversesTheEntireChain(t *testing.T) {
	e1 := io.EOF
	e2 := Wrap(e1, randString(), randomKeyValuePairs()...)
	e3 := Wrap(e2, randString(), randomKeyValuePairs()...)

	expected := e1
	actual := Unwrap(e3)

	if expected != actual {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func randString() string {
	return strconv.Itoa(rand.Intn(9999))
}

func randError() *StructuredError {
	return &StructuredError{
		msg:           randString(),
		cause:         nil,
		keyValuePairs: randomKeyValuePairs(),
	}
}

func randomKeyValuePairs() []interface{} {
	const l = 6
	kvs := make([]interface{}, l)
	for i := 0; i < l; i++ {
		kvs[i] = randString()
	}
	return kvs
}

func compareKVs(actual, expected []interface{}) bool {
	if len(actual) != len(expected) {
		return false
	}
	for _, v := range actual {
		found := false
		for _, v2 := range expected {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
