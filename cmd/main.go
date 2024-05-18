package main

import (
	"encoding/json"
	"flag"
	"fmt"

	errHelpers "MovieRecommenderCommandLine.IyadElwy/internal/errors"
	"MovieRecommenderCommandLine.IyadElwy/internal/rabbitmq"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	env, err := godotenv.Read()
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to read env variables")
	}

	genre := flag.String("genre", "", "A genre to search for. If not specified random genres will be used")
	numGenres := flag.Int("number_genres", 3, "Number of random genres (max 5 min 1)")
	flag.Parse()
	var genres []int
	if *genre == "" {
		if *numGenres <= 0 {
			panic("random genres: min 1")
		}
		randomGenres := returnNRandomGenres(env["GENRES_JSON_FILE_PATH"], *numGenres)
		genres = randomGenres
	} else {
		genreId, err := getIdByGenreName(env["GENRES_JSON_FILE_PATH"], *genre)
		if err != nil {
			panic(fmt.Sprintf("Genre %s not found", *genre))
		}
		genres = []int{genreId}
	}

	ch, err := rabbitmq.CreateChannel(env)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to open a channel")
	}

	// declare consume queue (get reply on)
	receiver_queue := env["RECEIVER_QUEUE"]
	q, err := createConsumeQueue(receiver_queue, ch)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to declare a queue")
	}

	// define response channel
	responsesChannel, err := ch.Consume(q.Name, receiver_queue, true, false, false, false, nil)
	if err != nil {
		errHelpers.FailedOnError(err, "Failed to register a consumer")
	}

	var movieDataList []MovieData
	for i := 0; i < *numGenres; i++ {
		workers_task_queue := env["WORKER_TASK_QUEUE"]
		err = sendTask(workers_task_queue, ch, q.Name, &rabbitmq.Task{
			GenreIds: []int{genres[i]},
		})
		if err != nil {
			errHelpers.FailedOnError(err, "Failed to publish a message")
		}

		movies, err := waitForResponse(responsesChannel)
		if err != nil {
			panic(err)
		}
		movieDataList = append(movieDataList, movies.Results...)
	}
	for _, m := range movieDataList {
		fmt.Println(m.Title)
	}

}

func createConsumeQueue(receiver_queue string, ch *amqp.Channel) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(receiver_queue, false, false, true, false, nil)
	if err != nil {
		return &q, err
	}
	return &q, nil
}

func sendTask(workers_task_queue string, ch *amqp.Channel, replyToQueueName string, task *rabbitmq.Task) error {
	err := ch.Publish("", workers_task_queue, false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: "1",
			Body:          task.ToJson(),
			ReplyTo:       replyToQueueName,
		})
	if err != nil {
		return err
	}
	return nil
}

func waitForResponse(responsesChannel <-chan amqp.Delivery) (MovieDataList, error) {
	for d := range responsesChannel {
		if d.CorrelationId == "1" {
			var movieDataList MovieDataList
			err := json.Unmarshal(d.Body, &movieDataList)
			if err != nil {
				return MovieDataList{}, err
			}
			return movieDataList, nil
		}
	}
	return MovieDataList{}, ErrNotFound
}
