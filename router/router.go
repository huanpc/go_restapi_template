package router

import (
	"net/http"
	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"

	"apistream/handler"
	"apistream/log"
	myMiddleware "apistream/router/middleware"
	"apistream/view"
)

func Register(r *chi.Mux) {

	// Setup log
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	// Add Middleware for router
	r.Use(chiMiddleware.Compress(2, "gzip"))
	r.Use(myMiddleware.CORS)

	r.Use(log.NewStructuredLogger(logger))

	//
	r.Group(func(r chi.Router) {
		r.Get("/apis/about", makeHandler(handler.About))
	})

	// API authentication
	r.Group(func(r chi.Router) {
		r.Post("/apis/auth", makeHandler(handler.Auth))
	})
}

func makeHandler(handlerFunc handler.HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := handler.GetBaseHandler(w, r)
		res := handlerFunc(h)
		view.RenderJson(w, res)
	}
}
