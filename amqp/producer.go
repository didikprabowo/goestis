package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/k/config"
	"github.com/k/utils"
	"log"
	"net/http"
	"time"
)

var co config.AllConfig
var urlConnection string

func init() {

	co = config.GetConfig()
	c := co.Worker
	urlConnection = c.URLConfig()

}

// SubcribeArticle
func SubcribeArticle(w http.ResponseWriter, r *http.Request, value interface{}) {

	data, _ := json.Marshal(value)

	conn, err := amqp.Dial(urlConnection)

	if err != nil {
		fmt.Println("Gagal Konek om dari produces...", urlConnection, err.Error())
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

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			DeliveryMode:    amqp.Persistent,
			Body:            data,
		})

	log.Printf(" [=>] Dikirim %s", time.Now())

	if err != nil {
		fmt.Println("Pesan tidak terkirim..")
	}

	out := map[string]interface{}{
		"results": data,
		"meta": map[string]interface{}{
			"per_page": 10,
			"status":   201,
			"message":  "ok",
		},
	}

	utils.RespondJSON(w, 201, out)
	return

}
