package repository

import (
	"context"
	"encoding/json"
	"github.com/k/src/database"
	"github.com/k/src/models"
	"strconv"
	"time"
)

// SaveMySQL
func SaveMySQL(payload []byte) {

	var news models.Article

	json.Unmarshal(payload, &news)

	db := database.GetConnection()

	stmtIns, err := db.Prepare("INSERT INTO news(author,body,created) VALUES( ?, ?,? )")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	res, err := stmtIns.Exec(news.Author, news.Body, news.Created)
	if err != nil {
		panic(err.Error())
	}

	lastID, _ := res.LastInsertId()
	SaveES(lastID, news.Created)

}

// SaveES
func SaveES(id int64, crt time.Time) {
	toES := struct {
		ID      int64     `json:"id"`
		Created time.Time `json:"created"`
	}{}

	toES.ID = id
	toES.Created = crt

	idIndex := strconv.FormatInt(id, 10)

	client := database.GetConnectionES()
	ctx := context.Background()

	_, err := client.Index().
		Index("news").
		Type("news").
		Id(idIndex).
		BodyJson(toES).
		Do(ctx)

	if err != nil {
		// Handle error
		panic(err)
	}

}
