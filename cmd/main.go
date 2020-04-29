package main

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"rabbit_write/internal/repository"
	"strings"
)

var (
	mode int
	message string
	filename string
)

func init() {
	flag.StringVar(&filename, "env", ".env", ".env file path")
	flag.StringVar(&message, "message", "test", "publish message")
	flag.IntVar(&mode, "mode", 1, "0 - publish, 1 - consume")
	flag.Parse()

	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	v := viper.New()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	client,err := repository.NewRabbitRepository(v)
	if err != nil {
		log.Fatal(err)
	}
	switch mode {
	case 0:
		err = client.Publish([]byte(message))
		if err != nil {
			log.Fatal(err)
		}
	case 1:
		err = client.Consume()
		if err != nil {
			log.Fatal(err)
		}
	}
}