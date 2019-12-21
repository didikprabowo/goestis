package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/k/amqp"
	"github.com/k/src/database"
	"github.com/k/src/models"
	"github.com/k/utils"
	"runtime"
	"sync"
	// "log"
	"net/http"
	"strconv"
	"time"
)

type (
	toES struct {
		ID      int64     `json:"id"`
		Created time.Time `json:"created"`
	}
	allNews struct {
		DataActicle []models.Article
	}
)

func AddItem(data models.Article, wg *sync.WaitGroup, news *allNews) {
	defer wg.Done()
	var mtx sync.Mutex

	mtx.Lock()
	local := append(news.DataActicle, data)
	news.DataActicle = local
	mtx.Unlock()
}

// PostArticle
func PostArticle(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		article := models.NewArticle()
		err := decoder.Decode(&article)

		amqp.SubcribeArticle(w, r, article)

		if err != nil {
			panic(err)
		}
	} else {
		runtime.GOMAXPROCS(2)
		keys, ok := r.URL.Query()["page"]

		var wg = &sync.WaitGroup{}
		var page string
		var skip int
		if !ok {
			skip = 0
			page = fmt.Sprintf("%v", 1)
		} else {
			skipInt, _ := strconv.Atoi(keys[0])
			skip = 10 * skipInt
			page = keys[0]
		}

		checkCache := getCache(page)

		var data allNews
		data = checkCache

		var code = 200

		if len(data.DataActicle) > 0 {
			fmt.Println("Data From Cache...")
			getCache := getCache(page)
			data = getCache
		} else {
			fmt.Println("Data From DB...")
			data, code = getEs(skip, wg)
			if len(data.DataActicle) > 0 {
				toCache, _ := json.Marshal(data.DataActicle)
				setCache(page, toCache)
			}

		}

		out := map[string]interface{}{
			"results": data.DataActicle,
			"meta": map[string]interface{}{
				"per_page": 10,
				"status":   code,
				"message":  "ok",
			},
		}

		utils.RespondJSON(w, code, out)
		return
	}
}

// producer
func producer(ch chan<- toES, docs toES) {
	ch <- docs

}

// getEs
func getEs(skip int, wg *sync.WaitGroup) (allNews, int) {
	fmt.Println("Search data in ES....")
	client := database.GetConnectionES()
	ctx := context.Background()

	var toApp allNews
	var statusCode int

	result, err := client.Search().
		Index("news").
		Sort("id", true).
		From(skip).Size(10).
		Pretty(true).
		Do(ctx)

	if err != nil {
		return toApp, http.StatusInternalServerError
	}

	if result.Hits.TotalHits.Value > 0 {
		var job = make(chan toES, 1)

		var results = make(chan models.Article, 1)

		go func() {
			for getData := range job {
				tostring := strconv.FormatInt(getData.ID, 10)
				toInt, err := strconv.Atoi(tostring)

				if err != nil {
					panic(err)
				}
				getNew(toInt, results, wg)
			}
		}()

		go func() {
			for resultz := range results {
				AddItem(resultz, wg, &toApp)
			}
		}()

		var docs toES

		for _, hit := range result.Hits.Hits {
			json.Unmarshal(hit.Source, &docs)
			producer(job, docs)
			println()
			wg.Add(2)
		}

	}
	wg.Wait()

	statusCode = http.StatusOK

	if len(toApp.DataActicle) < 1 {
		statusCode = http.StatusNoContent
	}

	return toApp, statusCode
}

// getNew
func getNew(ID int, ch chan<- models.Article, wg *sync.WaitGroup) {

	var news models.Article

	db := database.GetConnection()
	query := `SELECT id, author, body FROM news WHERE id = ? LIMIT 1;`

	st, err := db.Prepare(query)
	if err != nil {

		fmt.Println("error om querinya")
	}

	defer st.Close()

	r, err := st.Query(ID)
	if err != nil {
		fmt.Println("ID nya gak ada om")
	}

	defer r.Close()

	for r.Next() {
		err = r.Scan(&news.ID, &news.Author, &news.Body)
		if err != nil {
			panic(err)
		}
	}

	ch <- news

	defer wg.Done()

}

// getCache
func getCache(pages string) allNews {
	client := database.ConnectionRedis()
	val, err := client.Get(pages).Result()
	if err != nil {
		fmt.Println(err)
	}

	usr := allNews{}
	err = json.Unmarshal([]byte(val), &usr.DataActicle)
	return usr
}

// setCache
func setCache(pages string, st []byte) {
	client := database.ConnectionRedis()
	err := client.Set(pages, st, 0).Err()

	if err != nil {
		fmt.Println(err)
	}
}
