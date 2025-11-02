package web

import (
	"encoding/json"
	"net/http"
)

func JsonResponseWriter(w http.ResponseWriter, err *HttpError) error {
	w.WriteHeader(err.Code)
	return json.NewEncoder(w).Encode(err)
}
