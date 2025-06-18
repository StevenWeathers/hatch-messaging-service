package http

import "net/http"

type Config struct {
	ListenAddress string
}

// Endpoints represents the http service and its endpoints
type Endpoints struct {
	config Config
	router *http.ServeMux
}
