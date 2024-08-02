package server

import (
	"net/http"
	"thirdparty-service/database"
	"thirdparty-service/environment"

	"github.com/go-chi/chi"
)

func MountServer(cfg *environment.Config, mongodbStore database.MongoDBStore) *chi.Mux {
	router := chi.NewRouter()

	httpHandler := NewHTTPHandler(cfg, mongodbStore)

	// service check
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome our third party service"))
	})

	// Route POST /third-party/payments performs debit or credit transaction
	router.Post("/third-party/payments", httpHandler.PostPaymentHandler)

	// Route GET /third-party/payments/:reference retrieves transaction by reference
	router.Get("/third-party/payments/{reference}", httpHandler.GetPaymentHandler)

	return router
}
