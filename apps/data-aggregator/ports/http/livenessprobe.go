// The default listener for new recalculation requests
package http

import (
	"net/http"

	"github.com/CoreumFoundation/CoreDEX-API/utils/logger"
)

func NewListener() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			logger.Infof("Ping")
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})
	http.ListenAndServe(":8888", nil)
}
