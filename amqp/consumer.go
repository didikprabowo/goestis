package amqp

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/k/config"
	"github.com/k/src/repository"
	"log"
	"time"
)

func init() {

	co = config.GetConfig()
	c := co.Worker
	urlConnection = c.URLConfig()

}

// PubArticle
func PubArticle() {

	conn, err := amqp.Dial(urlConnection)

	if err != nil {
		fmt.Println("Gagal Konek om...", urlConnection, err.Error())
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		fmt.Println("gagal membuka Chanel...")
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"article", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		true,      // no-wait
		nil,       // arguments
	)

	if err != nil {
		fmt.Println("Deklarasi tidak benar..")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		fmt.Println("Pesan tidak terkirim..")
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			repository.SaveMySQL(d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	<-forever

}
