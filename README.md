# Lika Queue

Пакет включает в себя реализацию компонента для работы с очередями, а так же простые, но эффективные инструменты, которые
помогут просто и изящно построить многопоточные приложения.

# Базовый компонент

Базовый компонент очередей используется для хранения и управления несколькими брокерами сообщений, а так же является промежуточным
звеном между пользователем и конкретным брокером.

Это позволяет в свою очередь иметь возможность использовать разные типы очередей для разных целей, и при необходимости
заменять одну на другую, не задумываясь о деталях реализации.

```go

package main

import (
	"fmt"
	"github.com/liana-go/queue"
)

func main() {
	var queueComponent queue.QueueInterface
	var broker queue.Broker

	queueComponent = queue.New()
	broker = queue.NewMemoryBroker(10000)

	queueComponent.Add("main", broker)

	myMessage1 := "my string"
	myMessage2 := map[string]interface{} {
		"key1": "Vladimir",
		"key2": map[string]string {
			"name": "Dzhamshud",
			"surname": "Moskvich",
		},
	}

	_ = queueComponent.Publish("testQueue", myMessage1, nil)
	_ = queueComponent.Publish("testQueue", myMessage2, nil)

	for {
		message, _ := queueComponent.Consume("testQueue", nil)

		if message != nil {
			fmt.Println(message.Data())
		} else {
			break
		}
	}
}


```


### Использование очередей в высококонкурентной среде
Очереди отлично подходят для работы в высококонкурентной среде, когда происходит параллельные процессы на чтение и запись.

```go

package main

import (
	"fmt"
	"github.com/liana-go/queue"
	"time"
)

func main() {
	queueComponent := queue.New()
	queueComponent.Add("mem", queue.NewMemoryBroker(10000))

	go publishMessages(queueComponent)

	go consumeMessages(queueComponent, 1)
	go consumeMessages(queueComponent, 2)

	go publishMessages(queueComponent)

	time.Sleep(3 * time.Second)
}

func publishMessages(queueComponent queue.QueueInterface) {
	i := 0

	for i < 1000 {
		_ = queueComponent.Publish("testQueue", i, nil)
		i++
	}
}

func consumeMessages(queueComponent queue.QueueInterface, thread int) {
	for {
		message, _ := queueComponent.Consume("testQueue", nil)

		if message != nil {
			fmt.Println(fmt.Sprintf("%s %d", message.Data(), thread))
		} else {
			break
		}

	}
}

```

# Worker

Зачастую использование очередей предполагает под собой чтение данных из нескольких параллельных потоков и выполнением идентичных
операций в каждом из них, для этого в пакете есть готовые инструменты для упрощения данной задачи.

В пакете есть `QueueWorker` который помогает читать данные из очереди и одновременно работает с базовым воркером из пакета
[threading](https://github.com/liana-go/threading)

[подробнее о QueueWorker](queue_worker.md)

```go
package main

import (
	"fmt"
	"github.com/liana-go/queue"
)

func main() {
	queueComponent := queue.New()
	queueComponent.Add("mem", queue.NewMemoryBroker(10000))

	go publishMessages(queueComponent)

	var worker queue.QueueWorker

	worker = queue.QueueWorker {
		ThreadCount: 5,
		IsInfinite: true,
		QueueName: "testQueue",
		Broker: queueComponent, // or you can use any broker by calling queueComponent.Broker("mySpecialBroker")
		Callable: func(message queue.MessageData)  {
			fmt.Println(message.Data())
		},
	}

	worker.Run()
}

func publishMessages(queueComponent queue.QueueInterface) {
	i := 0

	for i < 1000 {
		_ = queueComponent.Publish("testQueue", i, nil)
		i++
	}
}

```