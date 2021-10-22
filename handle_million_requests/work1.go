package handle_million_requests

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type W1 struct {
	WgSend       *sync.WaitGroup
	Wg           *sync.WaitGroup
	MaxNum       int
	Ch           chan string
	DispatchStop chan struct{}
}

func (w *W1) Dispatch(job string) {
	w.WgSend.Add(10 * w.MaxNum)
	for i := 0; i < 10*w.MaxNum; i++ {
		go func(i int) {
			defer w.WgSend.Done()

			select {
			case w.Ch <- fmt.Sprintf("%d", i):
				return
			case <-w.DispatchStop:
				logrus.Debugln("退出发送 job: ", fmt.Sprintf("%d", i))
				return
			}
		}(i)
	}
}

func (w *W1) StartPool() {
	if w.Ch == nil {
		w.Ch = make(chan string, w.MaxNum)
	}

	if w.WgSend == nil {
		w.WgSend = &sync.WaitGroup{}
	}

	if w.Wg == nil {
		w.Wg = &sync.WaitGroup{}
	}

	if w.DispatchStop == nil {
		w.DispatchStop = make(chan struct{})
	}

	w.Wg.Add(w.MaxNum)
	for i := 0; i < w.MaxNum; i++ {
		go func() {
			defer w.Wg.Done()
			for v := range w.Ch {
				logrus.Debugf("完成工作: %s \n", v)
			}
		}()
	}
}

func (w *W1) Stop() {
	close(w.DispatchStop)
	w.WgSend.Wait()

	close(w.Ch)
	w.Wg.Wait()
}

func DealW1(max int) {
	w := NewWorker(w1, max)
	w.StartPool()
	w.Dispatch("")

	w.Stop()
}
