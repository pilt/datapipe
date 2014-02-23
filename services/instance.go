package services

import (
	"time"
	"sync"
)

func MemOpLatency() time.Duration {
	return time.Duration(0)
}

func MemByteLatency() time.Duration {
	return 10 * time.Microsecond
}

func MemBytesLatency(bytes int) time.Duration {
	return time.Duration(bytes) * MemByteLatency()
}

func DiskOpLatency() time.Duration {
	return 10 * time.Millisecond
}

func DiskByteLatency() time.Duration {
	return time.Millisecond
}

func DiskBytesLatency(bytes int) time.Duration {
	return time.Duration(bytes) * DiskByteLatency()
}

type Instance struct {
	services *Services
	transform InstanceFunc
	Config map[string]string
	Disk *InstanceStorage
	Mem *InstanceStorage
	sync.Mutex
}

type InstanceStorage struct {
	services *Services
	data map[string]string
	capacity int
	free int
	opLatency time.Duration
	byteLatency time.Duration
	sync.Mutex
}

func (self *InstanceStorage) Read(key string) (s string, ok bool) {
	self.Lock()
	defer self.Unlock()

	s, ok = self.data[key]
	self.sleep(s)

	return
}

func (self *InstanceStorage) Write(key, value string) bool {
	self.Lock()
	defer self.Unlock()

	newFree := self.free - len(value)
	self.sleep(value)
	if newFree < 0 {
		return false
	} else {
		self.data[key] = value
		self.free = newFree
		return true
	}
}

func (self *InstanceStorage) Remove(key string) {
	self.Lock()
	defer self.Unlock()

	s, ok := self.data[key]
	if ok {
		delete(self.data, key)
		self.free = self.free + len(s)
		self.sleep(s)
	} else {
		self.sleepOp()
	}
}

func (self *InstanceStorage) Free() int {
	self.Lock()
	defer self.Unlock()

	self.sleepOp()
	return self.free
}

func NewMemStorage(services *Services) *InstanceStorage {
	return &InstanceStorage{
		services: services,
		data: make(map[string]string),
		capacity: services.InstanceMem,
		free: services.InstanceMem,
		opLatency: MemOpLatency(),
		byteLatency: MemByteLatency(),
	}
}

func NewDiskStorage(services *Services) *InstanceStorage {
	return &InstanceStorage{
		services: services,
		data: make(map[string]string),
		capacity: services.InstanceDisk,
		free: services.InstanceDisk,
		opLatency: DiskOpLatency(),
		byteLatency: DiskByteLatency(),
	}
}

func (self *InstanceStorage) sleep(send string) {
	bytes := len(send)
	if bytes == 0 {
		self.sleepOp()
	} else {
		sleep := self.opLatency + time.Duration(bytes) * self.byteLatency
		time.Sleep(sleep)
	}
}

func (self *InstanceStorage) sleepOp() {
	if self.opLatency != 0 {
		time.Sleep(self.opLatency)
	}
}

func newInstance(s *Services, t InstanceFunc, config map[string]string) *Instance {
	instance := &Instance{
		services: s,
		transform: t,
		Mem: NewMemStorage(s),
		Disk: NewDiskStorage(s),
		Config: config,
	}
	return instance
}

func (self *Instance) Call(in string) string {
	self.Lock()
	defer self.Unlock()
	return self.transform(self, in)
}
