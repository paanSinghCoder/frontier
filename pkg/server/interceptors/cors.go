package interceptors

import (
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
)

func WithCors(h http.Handler, origin string) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins(strings.Split(origin, ",")),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "DELETE"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "X-Requested-With", "Accept", "Accept-Language", "Content-Language", "Origin"}),
		handlers.AllowCredentials(),
	)(h)
}