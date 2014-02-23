package services

import (
	"sync"
	"time"
)

func QueueOpLatency() time.Duration {
	return 20 * time.Millisecond
}

func QueueByteLatency() time.Duration {
	return 2 * time.Millisecond
}

type Queue struct {
	nextId int
	messages map[int]*Message
	Visibility time.Duration
	sync.Mutex
}

type Message struct {
	Id int
	Body string
	Sent time.Time
	ReceiveCount int
	Visible time.Time
}

func newQueue(visibility time.Duration) *Queue {
	return &Queue{
		Visibility: visibility,
		messages: make(map[int]*Message),
		nextId: 1,
	}
}

func (self *Queue) Send(b string) {
	bytes := len(b)
	time.Sleep(QueueOpLatency() + time.Duration(bytes) * QueueByteLatency())

	self.Lock()
	now := time.Now()
	m := &Message{
		Id: self.nextId,
		Body: b,
		Sent: now,
		ReceiveCount: 0,
		Visible: now,
	}
	self.nextId++
	self.messages[m.Id] = m
	self.Unlock()
}

func (self *Queue) Get() (m *Message) {
	self.Lock()
	now := time.Now()
	m = nil
	for _, qm := range self.messages {
		if now.After(qm.Visible) {
			m = qm
			m.Visible = now.Add(self.Visibility)
			m.ReceiveCount++
			break
		}
	}
	self.Unlock()

	bytes := 0
	if m != nil {
		bytes = len(m.Body)
	}
	time.Sleep(QueueOpLatency() + time.Duration(bytes) * QueueByteLatency())
	return
}

func (self *Queue) Delete(m *Message) {
	time.Sleep(QueueOpLatency())

	self.Lock()
	delete(self.messages, m.Id)
	self.Unlock()
}

func (self *Queue) Length() (l int) {
	time.Sleep(QueueOpLatency())
	self.Lock()
	l = len(self.messages)
	self.Unlock()
	return
}
