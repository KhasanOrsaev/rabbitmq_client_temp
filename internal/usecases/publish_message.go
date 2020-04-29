package usecases

import (
	"github.com/streadway/amqp"
	"rabbit_write/internal/domain"
)

type RabbitInteractor struct {
	Rabbit domain.Rabbit
}

func (client *RabbitInteractor) Publish(s []byte) error {
	var exchange, key string
	if client.Rabbit.Connection == nil || client.Rabbit.Connection.IsClosed() {
		err := client.Rabbit.OpenConnection()
		if err != nil {
			return err
		}
	}
	if client.Rabbit.Configuration.Exchange != nil {
		err := client.Rabbit.Channel.ExchangeDeclare(
			client.Rabbit.Configuration.Exchange.Name,
			client.Rabbit.Configuration.Exchange.Type,
			client.Rabbit.Configuration.Exchange.Durable,
			client.Rabbit.Configuration.Exchange.AutoDeleted,
			client.Rabbit.Configuration.Exchange.Internal,
			client.Rabbit.Configuration.Exchange.NoWait,
			client.Rabbit.Configuration.Exchange.Arguments,
		)
		if err != nil {
			return err
		}
		exchange = client.Rabbit.Configuration.Exchange.Name
		key = client.Rabbit.Configuration.Exchange.RoutingKey
	} else if client.Rabbit.Configuration.Queue != nil {
		q, err := client.Rabbit.Channel.QueueDeclare(
			client.Rabbit.Configuration.Queue.Name,
			client.Rabbit.Configuration.Queue.Durable,
			client.Rabbit.Configuration.Queue.AutoDelete,
			client.Rabbit.Configuration.Queue.Exclusive,
			client.Rabbit.Configuration.Queue.NoWait,
			client.Rabbit.Configuration.Queue.Arguments,
		)
		if err!=nil {
			return err
		}
		key = q.Name
	}
	err := client.Rabbit.Channel.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			Body:         s,
		},
	)
	if err != nil {
		return err
	}
	return nil
}