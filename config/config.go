package config

import (
	"os"
	"sync"
)

var once sync.Once

// Configuration struct for the app
type Configuration struct {
	Enviroment    string
	DB            string
	HTTPPort      string
	MongoURI      string
	MongoDB       string
	RabbitMQUser  string
	RabbitMQPass  string
	RabbitMQHost  string
	RabbitMQPort  string
	SecurityToken string
}

var (
	instance Configuration
)

// Instance create configuration instance
func Instance() Configuration {

	once.Do(func() {
		instance = loadConfig()
	})

	return instance
}

func loadConfig() Configuration {
	return Configuration{
		Enviroment:    os.Getenv("UP_ENV"),
		DB:            os.Getenv("UP_DB"),
		HTTPPort:      os.Getenv("UP_HTTP_PORT"),
		MongoURI:      os.Getenv("UP_MONGO_URI"),
		MongoDB:       os.Getenv("UP_MONGO_DB"),
		RabbitMQUser:  os.Getenv("UP_RABBITMQ_USER"),
		RabbitMQPass:  os.Getenv("UP_RABBITMQ_PASS"),
		RabbitMQHost:  os.Getenv("UP_RABBITMQ_HOST"),
		RabbitMQPort:  os.Getenv("UP_RABBITMQ_PORT"),
		SecurityToken: os.Getenv("UP_SECURITY_SECRET"),
	}
}
