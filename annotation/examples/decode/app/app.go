package app

import (
	"encoding/json"
	"net/http"
)

// FooMethod Lorem ipsum dolor sit amet, consectetur adipiscing elit
// @RestResponse(code=201, type=github.com/xgodev/boost/inject/examples/simple.Response, description=tiam sed efficitur purus at lacinia magna)
func FooMethod(w http.ResponseWriter, r *http.Request) {

	resp := Response{Message: "Hello World"}
	jsonData, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "error on convert json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsonData); err != nil {
		http.Error(w, "error on write response", http.StatusInternalServerError)
	}
}

type Response struct {
	Message string
}
