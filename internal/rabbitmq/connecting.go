package rabbitmq

import (
	"fmt"

	errHelpers "MovieRecommenderCommandLine.IyadElwy/internal/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func createConnection(env map[string]string) *amqp.Connection {
	rabbitHost := env["RABBITMQ_HOST"]
	rabbitPort := env["RABBITMQ_PORT"]
	rabbitUser := env["RABBITMQ_DEFAULT_USER"]
	rabbitPassword := env["RABBITMQ_DEFAULT_PASS"]

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		rabbitUser, rabbitPassword, rabbitHost, rabbitPort))
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to connect to RabbitMQ")
	}
	return conn
}

func CreateChannel(env map[string]string) (*amqp.Channel, error) {
	conn := createConnection(env)
	ch, err := conn.Channel()
	if err != nil {
		return ch, err
	}
	return ch, nil
}
