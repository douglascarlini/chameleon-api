package queue

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func (pipe *Pipe) Publish(data any) (string, error) {

	var err error
	var id string
	var message []byte

	id = uuid.NewString()

	if message, err = json.Marshal(data); err != nil {
		return "", err
	}

	pub := amqp.Publishing{ContentType: "application/json", Body: message, MessageId: id}
	err = pipe.Channel.Publish("", pipe.Queue.Name, false, false, pub)

	if err != nil {
		return "", err
	}

	return id, nil

}
