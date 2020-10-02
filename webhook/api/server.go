package api

import (
	"fmt"
	"net/http"
	"os"
)

// StartServer starts the server
func StartServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := &App{}

	mux := BuildRouter(app)

	fmt.Printf("Listening on port %s\n", port)

	return http.ListenAndServeTLS(fmt.Sprintf(":%s", port), "/tls/tls.crt", "/tls/tls.key", mux)
}
