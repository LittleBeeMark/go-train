package handle_million_requests

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type SubWorker struct {
	JobChan chan string
}

func (sw *SubWorker) Run(wg *sync.WaitGroup, poolCh chan chan string, quitCh chan struct{}) {
	if sw.JobChan == nil {
		sw.JobChan = make(chan string)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			poolCh <- sw.JobChan

			select {
			case res := <-sw.JobChan:
				logrus.Debugf("完成工作: %s \n", res)

			case <-quitCh:
				logrus.Debugf("消费者结束...... \n")
				return

			}
		}
	}()
}

type W2 struct {
	SubWorkers []SubWorker
	Wg         *sync.WaitGroup
	MaxNum     int
	ChPool     chan chan string
	QuitChan   chan struct{}
}

func (w *W2) Dispatch(job string) {
	jobChan := <-w.ChPool

	select {
	case jobChan <- job:
		logrus.Debugf("发送任务 : %s 完成 \n", job)
		return

	case <-w.QuitChan:
		logrus.Debugf("发送者（%s）结束 \n", job)
		return

	}
}

func (w *W2) StartPool() {
	if w.ChPool == nil {
		w.ChPool = make(chan chan string, w.MaxNum)
	}

	if w.SubWorkers == nil {
		w.SubWorkers = make([]SubWorker, w.MaxNum)
	}

	if w.Wg == nil {
		w.Wg = &sync.WaitGroup{}
	}

	for i := 0; i < len(w.SubWorkers); i++ {
		w.SubWorkers[i].Run(w.Wg, w.ChPool, w.QuitChan)
	}
}

func (w *W2) Stop() {
	close(w.QuitChan)
	w.Wg.Wait()

	close(w.ChPool)
}

func DealW2(max int) {
	w := NewWorker(w2, max)
	w.StartPool()
	for i := 0; i < 10*max; i++ {
		go w.Dispatch(fmt.Sprintf("%d", i))
	}

	w.Stop()
}
