package usecases

import (
	"fmt"
	"sync/atomic"
)

func (client *RabbitInteractor) Consume() error {
	if client.Rabbit.Connection == nil || client.Rabbit.Connection.IsClosed() {
		err := client.Rabbit.OpenConnection()
		if err != nil {
			return err
		}

	}
	q, err := client.Rabbit.Channel.QueueDeclare(
		client.Rabbit.Configuration.Queue.Name,    // name
		client.Rabbit.Configuration.Queue.Durable, // durable
		client.Rabbit.Configuration.Queue.AutoDelete, // delete when usused
		client.Rabbit.Configuration.Queue.Exclusive,  // exclusive
		client.Rabbit.Configuration.Queue.NoWait, // no-wait
		client.Rabbit.Configuration.Queue.Arguments,   // arguments
	)
	if err != nil {
		return err
	}
	err = client.Rabbit.Channel.Qos(
		100,
		0,
		false,
	)
	if err != nil {
		return err
	}
	msgs, err := client.Rabbit.Channel.Consume(
		q.Name, // queue name
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	i := int32(0)
	if q.Messages == 0 {
		return err
	}

	for d := range msgs {
		atomic.AddInt32(&i,1)
		fmt.Println(string(d.Body))
		if q.Messages == 0 || int(atomic.LoadInt32(&i)) >= q.Messages {
			break
		}
	}
	return nil
}
