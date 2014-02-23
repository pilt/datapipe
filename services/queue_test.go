package services

import (
	"testing"
	"time"
)

func TestQueueLatency(t *testing.T) {
	body := "01234567890asdfgbhjkoqiweuhdoiuhef"
	vis := 10 * QueueOpLatency()
	q := newQueue(vis)

	op := BucketOpLatency()
	bytes := len(body)
	shouldBe := 3 * op + 2 * time.Duration(bytes) * QueueByteLatency()

	concurrent := 150
	tic := time.Now()
	done := make(chan int)

	// A bunch of simultaneous calls.
	for i := 0; i != concurrent; i++ {
		go func() {
			q.Send(body)
			q.Get()
			q.Length()
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

func TestQueue(t *testing.T) {
	vis := 10 * QueueOpLatency()
	q := newQueue(vis)

	if q.Get() != nil {
		t.Error("expected no message")
		return
	}

	q.Send("foo")

	if q.Length() != 1 {
		t.Error("bad queue length")
		return
	}

	m := q.Get()
	if m == nil {
		t.Error("expected message")
		return
	}

	if q.Length() != 1 {
		t.Error("bad queue length")
		return
	}

	if m.Id != 1 {
		t.Error("bad message id")
		return
	}

	if m.ReceiveCount != 1 {
		t.Error("bad receive count")
	}

	if q.Get() != nil {
		t.Error("expected no message")
	}

	time.Sleep(vis)

	m = q.Get()
	if m == nil {
		t.Error("expected message")
		return
	}

	if m.ReceiveCount != 2 {
		t.Error("bad receive count")
	}

	q.Delete(m)

	if q.Get() != nil {
		t.Error("expected no message")
		return
	}

	time.Sleep(vis)

	if q.Get() != nil {
		t.Error("expected no message")
		return
	}

	if q.Length() != 0 {
		t.Error("bad queue length")
		return
	}
}

func TestMessageIds(t *testing.T) {
	vis := 10 * QueueOpLatency()
	q := newQueue(vis)

	q.Send("foo")
	q.Send("bar")
	q.Send("baz")

	if m := q.Get(); m == nil || m.Id != 1 {
		t.Error("bad message 1")
		return
	}
	if m := q.Get(); m == nil || m.Id != 2 {
		t.Error("bad message 2")
		return
	}
	if m := q.Get(); m == nil || m.Id != 3 {
		t.Error("bad message 3")
		return
	}
	if m := q.Get(); m != nil {
		t.Error("should not get a message")
		return
	}
}
