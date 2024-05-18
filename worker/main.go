package main

import (
	"log"

	"MovieRecommenderCommandLine.IyadElwy/internal/rabbitmq"

	errHelpers "MovieRecommenderCommandLine.IyadElwy/internal/errors"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	env, err := godotenv.Read()
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to read env variables")
	}

	ch, err := rabbitmq.CreateChannel(env)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to open a channel")
	}
	err = ch.Qos(1, 0, false)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to set QoS")
	}

	q, err := createConsuemQueue(env["WORKER_TASK_QUEUE"], ch)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to declare a queue")
	}

	// define task consumer
	tasksChannel, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to register a consumer")
	}

	startWorkerLoop(env, 3, tasksChannel, ch)
}

func createConsuemQueue(workers_task_queue string, ch *amqp.Channel) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(workers_task_queue, false, false, false, false, nil)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func startWorkerLoop(env map[string]string, numConsumers int, tasksChannel <-chan amqp.Delivery, ch *amqp.Channel) {
	for i := 0; i < numConsumers; i++ {
		go func() {
			for d := range tasksChannel {
				movies, err := getMovies(env["BASE_API_URL"], env["API_KEY"], rabbitmq.JsonToTask(d.Body).GenreIds)
				err = ch.Publish(
					"",
					d.ReplyTo,
					false,
					false,
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(movies),
					})
				if err != nil {
					errHelpers.FailedOnError(err, "Failed to publish a messsage")
				}
				d.Ack(false)
			}
		}()
	}

	log.Printf("Workers Awaiting Tasks...")
	// keep running workers
	for {
	}
}
