package schedule

import (
	"time"
)

type schedule struct {
	funcHandle func()
	seconds    time.Duration
}

func NewSchedule(f func(), seconds time.Duration) *schedule {
	return &schedule{f, seconds}
}

func (s *schedule) run() {
	d := time.Duration(time.Second * s.seconds)
	t := time.NewTicker(d)
	defer t.Stop()
	flag := make(chan int, 1)
	for {
		flag <- 1
		s.funcHandle()
		<-flag
		<-t.C
	}
}
