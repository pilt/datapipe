package services

import (
	"testing"
	"time"
)

func TestStorageBuckets(t *testing.T) {
	s := New(10, 100)

	b1 := s.GetBucket("foo")
	b2 := s.GetBucket("foo")

	if b1 != b2 {
		t.Error("not same bucket")
	}
}

func TestStorageQueues(t *testing.T) {
	s := New(10, 100)

	vis := 500 * time.Millisecond
	q1 := s.GetQueue("foo", vis)
	q2 := s.GetQueue("foo", vis)

	if q1 != q2 {
		t.Error("not same queue")
	}
}

func TestStorageVisibility(t *testing.T) {
	s := New(10, 100)

	vis := 500 * time.Millisecond
	q := s.GetQueue("foo", vis)
	if q.Visibility != vis {
		t.Error("bad visibility")
	}
}

func TestStorageVisibilityPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("should panic on new visibility")
		}
	}()

	s := New(10, 100)
	s.GetQueue("foo", 1 * time.Millisecond)
	s.GetQueue("foo", 2 * time.Millisecond)
}
