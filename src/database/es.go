package database

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/k/config"
	"log"
	"os"
	"time"
)

var client *elastic.Client

// GetConnectionES
func GetConnectionES() *elastic.Client {

	co := config.GetConfig()
	c := co.ES
	urlConnectionES := c.URLConfig()

	client, err := elastic.NewClient(
		elastic.SetURL(urlConnectionES),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		fmt.Println("Gagal Konek es broo...", err.Error())

	}

	return client
}
