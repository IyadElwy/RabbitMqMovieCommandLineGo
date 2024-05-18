package rabbitmq

import "encoding/json"

type Task struct {
	GenreIds []int `json:"genre_ids"`
}

func (task *Task) ToJson() []byte {
	res, err := json.Marshal(task)
	if err != nil {
		panic(err)
	}
	return res
}

func JsonToTask(taskJson []byte) Task {
	var task Task
	err := json.Unmarshal(taskJson, &task)
	if err != nil {
		panic(err)
	}
	return task
}
