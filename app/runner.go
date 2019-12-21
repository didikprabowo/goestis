package app

import (
	"github.com/k/handlers"
	"net/http"
)

// RunningApp
func RunningApp() {

	http.HandleFunc("/news", handlers.PostArticle)

	server := new(http.Server)

	server.Addr = ":9090"

	http.ListenAndServe(":9090", nil)

}
