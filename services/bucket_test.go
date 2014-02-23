package services

import (
	"testing"
	"time"
	"strings"
)

func TestBucketLatencySetGetDelete(t *testing.T) {
	b := newBucket()

	key := "01234567890asdfgbhjkoqiweuhdoiuhef"
	vals := []string{"", strings.Repeat(key, 10)}

	shouldBe := time.Duration(0)
	op := BucketOpLatency()
	for i := range vals {
		val := vals[i]
		keyDur := time.Duration(len(key)) * BucketByteLatency()
		keyValDur := time.Duration(len(key) + len(val)) * BucketByteLatency()
		shouldBe = shouldBe + 2 * (op + keyDur) + 2 * (op + keyValDur)
	}

	concurrent := 150
	tic := time.Now()
	done := make(chan int)

	// A bunch of simultaneous calls.
	for i := 0; i != concurrent; i++ {
		go func() {
			for j := range vals {
				val := vals[j]
				b.KeyExists(key)
				b.Set(key, val)
				b.Get(key)
				b.Delete(key)
			}

			done<- i
		}()
	}
	for doneCount := 0; doneCount != concurrent; doneCount++ {
		<-done
	}

	was := time.Now().Sub(tic)
	if was < shouldBe {
		t.Error("latency too low")
	}
	if 9 * was > 10 * shouldBe {
		t.Error("latency too high ", was, " ", shouldBe)
	}
}

func TestBucketLatencyList(t *testing.T) {
	b := newBucket()

	repeats := []int{1, 2, 3}
	baseKey := "01234567890asdfgbhjkoqiweuhdoiuhef"
	bytes := 0

	for i := range repeats {
		key := strings.Repeat(baseKey, repeats[i])
		b.Set(key, "foo")
		bytes = bytes + len(key)
	}

	shouldBe := BucketOpLatency() + time.Duration(bytes) * BucketByteLatency()

	concurrent := 150
	tic := time.Now()
	done := make(chan int)

	// A bunch of simultaneous list keys calls.
	for i := 0; i != concurrent; i++ {
		go func() {
			b.Keys()
			done<- i
		}()
	}
	for doneCount := 0; doneCount != concurrent; doneCount++ {
		<-done
	}

	was := time.Now().Sub(tic)
	if was < shouldBe {
		t.Error("latency too low")
	}
	if 9 * was > 10 * shouldBe {
		t.Error("latency too high ", was, " ", shouldBe)
	}
}

func TestBucket(t *testing.T) {
	b := newBucket()

	if len(b.Keys()) != 0 {
		t.Error("keys found")
		return
	}

	b.Set("foo", "foo")

	if b.KeyExists("foo") == false {
		t.Error("key not found")
		return
	}

	if k, _ := b.Get("foo"); k != "foo" {
		t.Error("bad key value")
		return
	}

	if len(b.Keys()) != 1 {
		t.Error("bad key list result")
		return
	}

	b.Delete("foo")

	if b.KeyExists("foo") == true {
		t.Error("key found")
		return
	}

	if _, ok := b.Get("foo"); ok {
		t.Error("get should not be OK")
		return
	}
}

func TestBucketReplace(t *testing.T) {
	b := newBucket()

	b.Set("foo", "foo")
	b.Set("foo", "bar")

	if k, _ := b.Get("foo"); k != "bar" {
		t.Error("bad key value")
		return
	}
}

func TestBucketManyKeys(t *testing.T) {
	b := newBucket()

	b.Set("foo", "foo")
	b.Set("bar", "bar")

	if len(b.Keys()) != 2 {
		t.Error("bad key count")
		return
	}
}
