package broker

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
	Queue *amqp.Queue
}



func NewRabbitMQ(URL string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(URL)

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get channel")
	}

	q, err := ch.QueueDeclare(
		"main", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to declare a queue.")
	}

	return &RabbitMQ{Connection: conn, Channel: ch, Queue: &q}, nil
}

func (r *RabbitMQ) Close()  {
	r.Channel.Close()
	r.Connection.Close()
}

type Message map[string]interface{}

func MsgSerialize(msg Message) ([]byte, error) {
    var b bytes.Buffer
    encoder := json.NewEncoder(&b)
    err := encoder.Encode(msg)
    return b.Bytes(), err
}

func MsgDeserialize(b []byte) (Message, error) {
    var msg Message
    buf := bytes.NewBuffer(b)
    decoder := json.NewDecoder(buf)
    err := decoder.Decode(&msg)
    return msg, err
}

func (r *RabbitMQ) Publish(msg []byte) error {
	err := r.Channel.Publish(
		"",     // exchange
		r.Queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	if err != nil {
		return errors.Wrap(err, "failed to publish a message.")
	}

	return nil
}
