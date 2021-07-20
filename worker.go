package queue

import (
	"github.com/liana-go/threading"
	"time"
)

type QueueWorker struct {
	ThreadCount     int
	QueueName       string
	Duration        time.Duration // milliseconds. default 100 ms
	ConsumingParams map[string]interface{}
	Broker          Broker
	Callable        func(message MessageData)
	IsInfinite      bool // if set false, message consuming will stop after getting first nil result from broker
	worker          *threading.Worker
}

func (qw *QueueWorker) Run() {
	qw.initWorker()

	if qw.IsInfinite && qw.Duration <= 0 {
		qw.Duration = 100
	}
	qw.worker.Callable = qw.consume
	qw.worker.Run()
}

func (qw *QueueWorker) consume() {
	for {
		m, err := qw.Broker.Consume(qw.QueueName, qw.ConsumingParams)
		// todo add logger
		if err != nil {
			return
		}

		if m == nil {
			if qw.IsInfinite {
				time.Sleep(qw.Duration * time.Millisecond)
				continue
			} else {
				break
			}
		}

		qw.Callable(m)
	}
}

func (qw *QueueWorker) initWorker() {
	if qw.worker == nil {
		qw.worker = &threading.Worker{
			ThreadCount: 5,
			Callable: func() {
				println(true)
			},
		}
	}
}
