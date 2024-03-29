package port

import (
	"fmt"
	"sync"
)

type MyPortPool struct {
	mutex sync.Mutex

	freePortSet *Set[uint16] // store rtp ports (even number)
	start, end  uint16
}

func NewPortPool() *MyPortPool {
	return &MyPortPool{
		freePortSet: NewSet[uint16](),
	}
}

func (p *MyPortPool) Init(start uint16, end uint16) {
	// NOTE: rtp use even port
	if start&0x01 != 0 {
		start += 1
	}
	if end&0x01 != 0 {
		end -= 1
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for i := start; i < end; i += 2 {
		p.freePortSet.Add(i)
	}

	p.start = start
	p.end = end
}

func (p *MyPortPool) Get() (port uint16) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// if pool is empty, return 0 as this port wouldn't be used by applications

	if p.freePortSet.Size() == 0 {
		goto noPort
	}
	port = p.freePortSet.GetAndRemove()
noPort:
	return
}

func (p *MyPortPool) Put(port uint16) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.freePortSet.Contain(port) {
		fmt.Printf("put port: %v to pool but it is already in it", port)
		return
	}
	p.freePortSet.Add(port)
}
