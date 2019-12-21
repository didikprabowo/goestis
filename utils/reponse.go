package utils

import (
	"encoding/json"
	"net/http"
)

// RespondJSON
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// RespondError
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

// func responseAsJSON(w http.ResponseWriter, response interface{}, code int) {
// 	js, err := json.Marshal(response)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	_, err = w.Write(js)
// 	if err != nil {
// 		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
// 	}

// }
