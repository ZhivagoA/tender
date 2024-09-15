package handlers

import (
	"net/http"
)

type PingHandler struct{}

func (h *PingHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
