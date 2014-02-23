package services

import (
	"sync"
)

type Group struct {
	services *Services
	config map[string]string
	instances []*Instance
	counter int
	sync.Mutex
}

func newGroup(s *Services, count int, t InstanceFunc, config map[string]string) *Group {
	instances := make([]*Instance, count, count)
	for i := 0; i < count; i++ {
		instances[i] = newInstance(s, t, config)
	}

	return &Group{
		services: s,
		config: config,
		instances: instances,
		counter: 0,
	}
}

func (self *Group) CallRoundRobin(in string) string {
	res := make(chan string)
	instance := self.instances[self.counter % len(self.instances)]
	go func() {
		res<- instance.Call(in)
	}()

	self.Lock()
	self.counter++
	self.Unlock()

	return <-res
}
