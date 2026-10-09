package server

import (
	"net/http"

	"github.com/coderianx/lantern/internal/server/handlers"
	"github.com/go-chi/chi/v5"
)

func Run() {
	router := chi.NewRouter()

	router.Route("/api/v1", func(r chi.Router) {
		r.Post(
			"/auth/login", handlers.Login,
		)
	})

	http.ListenAndServe(":9090", router)
}
