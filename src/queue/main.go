package queue

import (
	"github.com/streadway/amqp"
)

type Pipe struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

type RMQ struct {
	Conn *amqp.Connection
	Pipe map[string]*Pipe
}

func Connect(dsn string) *RMQ {

	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err)
	}

	return &RMQ{Conn: conn, Pipe: make(map[string]*Pipe)}

}

func (rmq *RMQ) NewQueue(name string) error {

	ch, err := rmq.Conn.Channel()
	if err != nil {
		rmq.Conn.Close()
		return err
	}

	q, err := ch.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		rmq.Conn.Close()
		ch.Close()
		return err
	}

	rmq.Pipe[name] = &Pipe{Channel: ch, Queue: q}

	return nil

}
