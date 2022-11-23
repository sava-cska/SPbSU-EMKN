package event_queue

import (
	"encoding/json"
	"fmt"
	"github.com/sava-cska/SPbSU-EMKN/internal/app/models"
	"github.com/streadway/amqp"
)

type EventQueue interface {
	AddMessage(name string, message models.Message) error
}

type eventQueueImpl struct {
	Config     *Config
	Connection *amqp.Connection
}

func New(config *Config, rabbitLogin, rabbitPassword string) (EventQueue, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitLogin, rabbitPassword, config.Address, config.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	return &eventQueueImpl{
		Config:     config,
		Connection: conn,
	}, nil
}

func (eventQueueImpl eventQueueImpl) AddMessage(name string, message models.Message) error {
	amqpChannel, err := eventQueueImpl.Connection.Channel()
	if err != nil {
		return err
	}
	queue, err := amqpChannel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return err
	}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	return err
}
