package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type (
	RabbitConfig struct {
		Host     string
		User     string
		Password string
		Port     int
		VHost    string
	}

	MySQLConfig struct {
		Host     string
		User     string
		Password string
		DB       string
	}
	ESConfig struct {
		Host string
		Port int
	}
	RedisConfig struct {
		Host     string
		Password string
		Port     string
		DB       string
	}

	AllConfig struct {
		Worker RabbitConfig
		DB     MySQLConfig
		ES     ESConfig
		Redis  RedisConfig
	}
)

var configApp AllConfig

func (co RabbitConfig) URLConfig() string {

	auth := fmt.Sprintf("%s:%s@", co.User, co.Password)
	return fmt.Sprintf("amqp://%s%s:%d/%s", auth, co.Host, co.Port, co.VHost)
}

func (co MySQLConfig) URLConfig() string {
	return fmt.Sprintf("%v:%v@/%v", co.User, co.Password, co.DB)
}

func (co ESConfig) URLConfig() string {
	return fmt.Sprintf("http://%s:%d", co.Host, co.Port)
}

func (co RedisConfig) URLConfig() string {
	return fmt.Sprintf("Addr:%v:%v, Password:%v,DB:%v", co.Host, co.Port, co.Password, co.DB)
}

func GetConfig() AllConfig {
	NewConfig()
	return configApp
}

func NewConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rabbitHost := os.Getenv("RABBIT_HOST")
	rabbitUser := os.Getenv("RABBIT_USER")
	rabbitPassword := os.Getenv("RABBIT_PASSWORD")
	rabbitPort := os.Getenv("RABBIT_PORT")
	rabbitVHost := os.Getenv("RABBIT_VHOST")

	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDB := os.Getenv("MYSQL_DB")

	esHost := os.Getenv("ES_HOST")
	esPort := os.Getenv("ES_PORT")

	reHost := os.Getenv("REDIS_HOST")
	rePass := os.Getenv("REDIS_PASSWORD")
	rePort := os.Getenv("REDIS_PORT")
	reDB := os.Getenv("REDIS_DB")

	//Rabit
	rb := RabbitConfig{}
	rb.Host = rabbitHost
	rb.User = rabbitUser
	rb.Password = rabbitPassword
	rb.Port, _ = strconv.Atoi(rabbitPort)
	rb.VHost = rabbitVHost

	//MySQL
	my := MySQLConfig{}

	my.Host = mysqlHost
	my.User = mysqlUser
	my.Password = mysqlPassword
	my.DB = mysqlDB

	// ES
	es := ESConfig{}
	es.Host = esHost
	es.Port, _ = strconv.Atoi(esPort)

	// Redis

	re := RedisConfig{}
	re.Host = reHost
	re.Password = rePass
	re.Port = rePort
	re.DB = reDB

	// Expose
	configApp.Worker = rb
	configApp.DB = my
	configApp.ES = es
	configApp.Redis = re

}

func init() {
	NewConfig()
}
