package handlers

import (
	"net/http"

	"github.com/Ali-Karaki/e8markets/server/internal/httpx"
)

func Health(w http.ResponseWriter, r *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
