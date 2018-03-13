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

	// Template
	r.Group(func(r chi.Router) {
		r.Get("/apis/about", makeHandler(handler.About))
	})

	// API which receive requests from nginx-rtmp
	r.Group(func(r chi.Router) {
		r.Post("/apis/auth", makeHandler(handler.Auth))
	})

	// API logging & tracking
	r.Group(func(r chi.Router) {
		r.Post("/apis/event", makeHandler(handler.Event))
	})
}

func makeHandler(handlerFunc handler.HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := handler.GetBaseHandler(w, r)
		res := handlerFunc(h)
		view.RenderJson(w, res)
	}
}
