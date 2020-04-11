package engine

import (
	"encoding/json"
	"net/http"
)

type errRes struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

// Respond is a wrapper to respond in json format easier
func Respond(w http.ResponseWriter, _ *http.Request, status int, data interface{}) {
	if e, ok := data.(error); ok {
		tmp := new(errRes)
		tmp.Error = e.Error()
		tmp.Status = http.StatusText(status)
		data = tmp
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}
