package services

import (
	"testing"
	"time"
	"strings"
)

func TestInstanceMemLatency(t *testing.T) {
	s := New(1e3, 1e4)
	m := NewMemStorage(s)
	send := strings.Repeat("01234567890asdfgbhjkoqiweuhdoiuhef", 100)

	tic := time.Now()
	shouldBe := MemOpLatency() + MemBytesLatency(len(send))
	m.Write("foo", send)
	was := time.Now().Sub(tic)
	if was < shouldBe {
		t.Error("latency too low")
	}
	if 9 * was > 10 * shouldBe {
		t.Error("latency too high ", was, " ", shouldBe)
	}
}

func TestInstanceDiskLatency(t *testing.T) {
	s := New(1e3, 1e4)
	d := NewDiskStorage(s)
	send := "01234567890asdfgbhjkoqiweuhdoiuhef"

	tic := time.Now()
	shouldBe := DiskOpLatency() + DiskBytesLatency(len(send))
	d.Write("foo", send)
	was := time.Now().Sub(tic)
	if was < shouldBe {
		t.Error("latency too low")
	}
	if 9 * was > 10 * shouldBe {
		t.Error("latency too high ", was, " ", shouldBe)
	}
}

func TestInstanceDiskLatencyConcurrent(t *testing.T) {
	s := New(1e3, 1e4)
	d := NewDiskStorage(s)
	send := "01234567890asdfgbhjkoqiweuhdoiuhef"
	concurrent := 5

	shouldBe := (time.Duration(concurrent) * 2 *
	             (DiskOpLatency() + DiskBytesLatency(len(send))))

	tic := time.Now()
	done := make(chan int)

	// Start routines that will attempt to read and write data concurrently.
	for i := 0; i != concurrent; i++ {
		go func() {
			d.Write("foo", send)
			d.Read("foo")
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

func TestInstanceDiskLatencyDelete(t *testing.T) {
	s := New(1e3, 1e4)
	d := NewDiskStorage(s)

	iters := 50
	shouldBe := time.Duration(iters) * DiskOpLatency()
	tic := time.Now()
	for i := 0; i != iters; i++ {
		d.Remove("a")
	}
	was := time.Now().Sub(tic)

	if was < shouldBe {
		t.Error("latency too low")
	}
	if 9 * was > 10 * shouldBe {
		t.Error("latency too high ", was, " ", shouldBe)
	}
}

func TestInstanceMemEmpty(t *testing.T) {
	s := New(5, 100)
	m := NewMemStorage(s)
	m.Write("foo", "012345")

	if m.Free() != 5 {
		t.Error("bad free memory")
		return
	}
}

func TestInstanceMemBasic(t *testing.T) {
	s := New(10, 100)
	m := NewMemStorage(s)

	if m.Free() != 10 {
		t.Error("bad initial free memory")
		return
	}

	m.Write("foo", "foo")

	if m.Free() != 10 - len("foo") {
		t.Error("bad free memory after write")
		return
	}

	m.Remove("foo")

	if m.Free() != 10 {
		t.Error("bad free memory after remove")
		return
	}
}
