package services

import (
	"sync"
	"time"
)

type InstanceFunc func(*Instance, string) string

type Services struct {
	buckets map[string]*Bucket
	queues map[string]*Queue
	InstanceMem int
	InstanceDisk int
	sync.Mutex
}

func New(memory, disk int) *Services {
	return &Services{
		buckets: make(map[string]*Bucket),
		queues:  make(map[string]*Queue),
		InstanceMem: memory,
		InstanceDisk: disk,
	}
}

// Compute

func (self *Services) NewInstance(t InstanceFunc, c map[string]string) *Instance {
	return newInstance(self, t, c)
}

func (self *Services) NewGroup(count int, t InstanceFunc, c map[string]string) *Group {
	return newGroup(self, count, t, c)
}

// Storage

func (self *Services) GetBucket(name string) (b *Bucket) {
	self.Lock()
	defer self.Unlock()

	b = self.buckets[name]
	if b == nil {
		self.buckets[name] = newBucket()
		b = self.buckets[name]
	}
	return b
}

func (self *Services) GetQueue(name string, visibility time.Duration) (q *Queue) {
	self.Lock()
	defer self.Unlock()

	q = self.queues[name]
	if q == nil {
		self.queues[name] = newQueue(visibility)
		q = self.queues[name]
	}
	if q.Visibility != visibility {
		panic("tried to get existing queue with new visibility")
	}
	return q
}
