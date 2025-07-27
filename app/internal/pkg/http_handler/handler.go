package http_handler

import (
	"encoding/json"
	"net/http"
	"simple_lgtm/internal/model"
	"simple_lgtm/internal/pkg/errs"
)

func AbortJSON(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}
	status, message := errs.MapHttp(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(model.Response{Message: message})
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(data)
}
