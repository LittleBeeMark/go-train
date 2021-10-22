package handle_million_requests

import (
	"sync"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
}

type Work interface {
	Dispatch(job string)
	StartPool()
	Stop()
}

const (
	w1 = "w1"
	w2 = "w2"
)

// NewWorker doc
func NewWorker(t string, max int) Work {
	switch t {
	case w1:
		return &W1{
			MaxNum: max,
			Wg:     &sync.WaitGroup{},
			Ch:     make(chan string),
		}

	case w2:
		return &W2{
			Wg:       &sync.WaitGroup{},
			MaxNum:   max,
			ChPool:   make(chan chan string, max),
			QuitChan: make(chan struct{}),
		}

	}

	return nil

}

