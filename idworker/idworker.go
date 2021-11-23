package idworker

import (
	"time"
)

const alf = 779639820323879

type Idworker struct {
	stock chan int64
}

func NewIdworker(len int64) *Idworker {
	Idworker := new(Idworker)
	Idworker.stock = make(chan int64, len)
	go Idworker.init()
	return Idworker
}

func (Idworker *Idworker) GetId() int64 {
	return <-Idworker.stock
}

func (Idworker *Idworker) init() {
	d := time.Duration(time.Millisecond)
	t := time.NewTicker(d)
	defer t.Stop()
	for {
		<-t.C
		Idworker.stock <- time.Now().UnixNano()/1e3 + alf
	}
}
