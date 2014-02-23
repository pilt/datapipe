package services

import (
	"sync"
	"time"
)

func BucketOpLatency() time.Duration {
	return 20 * time.Millisecond
}

func BucketByteLatency() time.Duration {
	return 2 * time.Millisecond
}

type Bucket struct {
	blobs map[string]string
	sync.Mutex
}

func newBucket() *Bucket {
	return &Bucket{
		blobs: make(map[string]string),
	}
}

func (self *Bucket) Keys() []string {
	time.Sleep(BucketOpLatency())

	self.Lock()
	keys := make([]string, len(self.blobs))
	i := 0
	bytes := 0
	for k, _ := range self.blobs {
		keys[i] = k
		bytes = bytes + len(k)
	}
	self.Unlock()

	time.Sleep(time.Duration(bytes) * BucketByteLatency())

	return keys
}

func (self *Bucket) KeyExists(key string) bool {
	bytes := len(key)
	time.Sleep(BucketOpLatency() + time.Duration(bytes) * BucketByteLatency())

	self.Lock()
	_, exists := self.blobs[key]
	self.Unlock()

	return exists
}

func (self *Bucket) Get(key string) (s string, ok bool) {
	inBytes := len(key)
	time.Sleep(BucketOpLatency() + time.Duration(inBytes) * BucketByteLatency())

	self.Lock()
	s, ok = self.blobs[key]
	self.Unlock()

	outBytes := len(s)
	time.Sleep(time.Duration(outBytes) * BucketByteLatency())
	return
}

func (self *Bucket) Set(key, value string) {
	bytes := len(key) + len(value)
	sleep := BucketOpLatency() + time.Duration(bytes) * BucketByteLatency()
	time.Sleep(sleep)

	self.Lock()
	self.blobs[key] = value
	self.Unlock()
}

func (self *Bucket) Delete(key string) {
	bytes := len(key)
	sleep := BucketOpLatency() + time.Duration(bytes) * BucketByteLatency()
	time.Sleep(sleep)

	self.Lock()
	delete(self.blobs, key)
	self.Unlock()
}
