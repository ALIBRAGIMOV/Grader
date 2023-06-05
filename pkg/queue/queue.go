package queue

import (
	"flag"
)

const (
	SolutionQueueName = "solution"
	ResultQueueName   = "result_solution"
)

var (
	RabbitAddr = flag.String("rabbit", "amqp://guest:guest@localhost:5672/", "rabbit addr")
)
