package domain

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type Rabbit struct {
	Configuration Configuration
	Connection    *amqp.Connection
	Channel       *amqp.Channel
}

type Configuration struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Tlc bool `json:"tlc"`
	VirtualHost string `json:"virtual_host"`
	Exchange *Exchange `json:"exchange"`
	Queue *Queue `json:"queue"`
	Workers int `json:"workers"`
}
type Exchange struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Durable bool `json:"durable"`
	AutoDeleted bool `json:"auto_deleted"`
	Internal bool `json:"internal"`
	NoWait bool `json:"no_wait"`
	Arguments map[string]interface{} `json:"arguments"`
	RoutingKey string `json:"routing_key"`
}

type Queue struct {
	Name string `json:"name"`
	Durable bool `json:"durable"`
	AutoDelete bool `json:"auto_delete"`
	Exclusive bool `json:"exclusive"`
	NoWait bool `json:"no_wait"`
	Arguments map[string]interface{} `json:"arguments"`
}

func NewRabbit(v *viper.Viper) (*Rabbit,error) {
	var conf Configuration
	v.SetDefault("rabbitmq.host","localhost")
	v.BindEnv("rabbitmq.host", "RABBITMQ_HOST")
	v.SetDefault("rabbitmq.port","5672")
	v.BindEnv("rabbitmq.port", "RABBITMQ_PORT")
	v.SetDefault("rabbitmq.user","csstat")
	v.BindEnv("rabbitmq.user", "RABBITMQ_USER")
	v.SetDefault("rabbitmq.password","Qweasd123")
	v.BindEnv("rabbitmq.password", "RABBITMQ_PASSWORD")
	v.SetDefault("rabbitmq.virtual.host","")
	v.BindEnv("rabbitmq.virtual.host", "RABBITMQ_VIRTUAL_HOST")
	v.SetDefault("rabbitmq.workers",1)
	v.BindEnv("rabbitmq.workers", "RABBITMQ_WORKERS")
	v.SetDefault("rabbitmq.tls",false)
	v.BindEnv("rabbitmq.tls", "RABBITMQ_TLS")
	// bind queue vars
	v.BindEnv("rabbitmq.queue.name")
	v.BindEnv("rabbitmq.queue.durable")
	v.BindEnv("rabbitmq.queue.auto.delete")
	v.BindEnv("rabbitmq.queue.exclusive")
	v.BindEnv("rabbitmq.queue.nowait")
	v.BindEnv("rabbitmq.queue.arguments")
	// bind exchange vars
	v.BindEnv("rabbitmq.exchange.name")
	v.BindEnv("rabbitmq.exchange.durable")
	v.BindEnv("rabbitmq.exchange.type")
	v.BindEnv("rabbitmq.exchange.routing.key")
	v.BindEnv("rabbitmq.exchange.auto.deleted")
	v.BindEnv("rabbitmq.exchange.internal")
	v.BindEnv("rabbitmq.exchange.nowait")
	v.BindEnv("rabbitmq.exchange.arguments")

	conf.Host = v.GetString("rabbitmq.host")
	conf.Port = v.GetString("rabbitmq.port")
	conf.User = v.GetString("rabbitmq.user")
	conf.Password = v.GetString("rabbitmq.password")
	conf.VirtualHost = v.GetString("rabbitmq.virtual.host")
	conf.Workers = v.GetInt("rabbitmq.workers")
	conf.Tlc = v.GetBool("rabbitmq.tls")


	if name := v.Get("rabbitmq.queue.name"); name != nil {
		conf.Queue = &Queue{}
		conf.Queue.Name = name.(string)
		if durable := v.Get("rabbitmq.queue.durable"); durable != nil {
			conf.Queue.Durable = v.GetBool("rabbitmq.queue.durable")
		}
		if autoDelete := v.Get("rabbitmq.queue.auto.delete"); autoDelete != nil {
			conf.Queue.AutoDelete = v.GetBool("rabbitmq.queue.auto.delete")
		}
		if exclusive := v.Get("rabbitmq.queue.exclusive"); exclusive != nil {
			conf.Queue.Exclusive = v.GetBool("rabbitmq.queue.exclusive")
		}
		if noWait := v.Get("rabbitmq.queue.nowait"); noWait != nil {
			conf.Queue.NoWait = v.GetBool("rabbitmq.queue.nowait")
		}
		if arguments := v.Get("rabbitmq.queue.arguments"); arguments != nil {
			conf.Queue.Arguments = v.GetStringMap("rabbitmq.queue.durable")
		}
	}

	if name := v.Get("rabbitmq.exchange.name"); name != nil {
		conf.Exchange = &Exchange{}
		conf.Exchange.Name = name.(string)
		if durable := v.Get("rabbitmq.exchange.durable"); durable != nil {
			conf.Exchange.Durable = v.GetBool("rabbitmq.exchange.durable")
		}
		if t := v.Get("rabbitmq.exchange.type"); t != nil {
			conf.Exchange.Type = t.(string)
		}
		if key := v.Get("rabbitmq.exchange.routing.key"); key != nil {
			conf.Exchange.RoutingKey = key.(string)
		}
		if autoDelete := v.Get("rabbitmq.exchange.auto.deleted"); autoDelete != nil {
			conf.Exchange.AutoDeleted = v.GetBool("rabbitmq.exchange.auto.deleted")
		}
		if internal := v.Get("rabbitmq.exchange.internal"); internal != nil {
			conf.Exchange.Internal = v.GetBool("rabbitmq.exchange.internal")
		}
		if noWait := v.Get("rabbitmq.exchange.nowait"); noWait != nil {
			conf.Exchange.NoWait = v.GetBool("rabbitmq.exchange.nowait")
		}
		if arguments := v.Get("rabbitmq.exchange.arguments"); arguments != nil {
			conf.Exchange.Arguments = v.GetStringMap("rabbitmq.exchange.arguments")
		}
	}

	return &Rabbit{Configuration: conf}, nil
}

func (client *Rabbit) OpenConnection() error {
	var err error
	if client.Configuration.Tlc {
		cfg := tls.Config{InsecureSkipVerify:true}
		client.Connection, err = amqp.DialTLS("amqps://" + client.Configuration.User + ":" + client.Configuration.Password +
			"@" + client.Configuration.Host + ":" + client.Configuration.Port + "/" + client.Configuration.VirtualHost, &cfg)
		if err != nil {
			return err
		}
		client.Channel, err = client.Connection.Channel()
		if err != nil {
			return err
		}
	} else {
		client.Connection, err = amqp.Dial("amqp://" + client.Configuration.User + ":" + client.Configuration.Password +
			"@" + client.Configuration.Host + ":" + client.Configuration.Port)
		if err != nil {
			return err
		}
		client.Channel, err = client.Connection.Channel()
		if err != nil {
			return err
		}
	}
	return nil
}

func (client *Rabbit) CloseConnection() error {
	var err error
	err = client.Channel.Close()
	if err != nil {
		return err
	}
	err = client.Connection.Close()
	if err != nil {
		return err
	}
	return nil
}