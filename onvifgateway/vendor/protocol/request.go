package protocol

import (
	"fmt"
	"sync"
)

type Request struct {
	locker sync.Mutex
	store  map[string]chan *OnvifFindInfo
}

func (this *Request) Init() {
	this.store = make(map[string]chan *OnvifFindInfo)
}

func (this *Request) NewWaiter(id string) (chan *OnvifFindInfo, error) {
	this.locker.Lock()
	defer this.locker.Unlock()

	// 值存在
	if _, ok := this.store[id]; ok {
		return nil, fmt.Errorf("value have exists")
	}

	// 申名新的对像
	c := make(chan *OnvifFindInfo)
	// 存储
	this.store[id] = c
	return c, nil
}

func (this *Request) CloseWaiter(id string, c chan *OnvifFindInfo) {
	this.locker.Lock()
	defer this.locker.Unlock()

	close(c)
	// 清理存储
	delete(this.store, id)
}

func (this *Request) Notify(id string, info *OnvifFindInfo) {
	this.locker.Lock()
	defer this.locker.Unlock()

	c, o := this.store[id]
	if false == o {
		return
	}

	c <- info
}
