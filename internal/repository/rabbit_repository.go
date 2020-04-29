package repository

import (
	"github.com/spf13/viper"
	"rabbit_write/internal/domain"
	"rabbit_write/internal/usecases"
)

type RabbitRepository interface {
	Publish(s []byte) error
	Consume() error
}

func NewRabbitRepository(v *viper.Viper) (RabbitRepository,error) {
	rabbit,err := domain.NewRabbit(v)
	if err != nil {
		return nil, err
	}
	return &usecases.RabbitInteractor{
		Rabbit: *rabbit,
	}, nil
}