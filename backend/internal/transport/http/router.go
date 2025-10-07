package http

import "net/http"

func SetupRoutes(handler *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/receipts/upload", handler.UploadAndProcess)
	mux.HandleFunc("/api/health", handler.HealthCheck)

	return enableCORS(mux)
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
