package main

import (
	"flag"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"grader/pkg/queue"
	queueDelivery "grader/pkg/queue/delivery"
	"grader/pkg/utils"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel
)

func main() {
	flag.Parse()
	var err error

	httpClient := &http.Client{
		Timeout: time.Duration(2 * time.Minute),
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	rabbitConn, err = amqp.Dial(*queue.RabbitAddr)
	utils.FatalOnError("cant connect to rabbit", err)

	rabbitChan, err = rabbitConn.Channel()
	utils.FatalOnError("cant open chan", err)
	defer rabbitChan.Close()

	_, err = rabbitChan.QueueDeclare(
		queue.SolutionQueueName, // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	utils.FatalOnError("cant init queue", err)

	err = rabbitChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FatalOnError("cant set QoS", err)

	solutions, err := rabbitChan.Consume(
		queue.SolutionQueueName, // queue
		"",                      // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	utils.FatalOnError("cant register consumer", err)

	queueHandler := &queueDelivery.QueueHandler{
		Client: httpClient,
		Logger: logger,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	numWorkers := runtime.NumCPU()

	for i := 0; i < numWorkers; i++ {
		go queueHandler.SolutionWorker(solutions)
	}

	log.Println("queue processor started")
	wg.Wait()
}
