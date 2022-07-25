package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

type server struct {
	router *chi.Mux
}

func New() *server {
	r := chi.NewRouter()

	return &server{router: r}
}

func (s *server) AddHandler(method, path string, handlerFn http.HandlerFunc) {
	s.router.MethodFunc(method, path, handlerFn)
}

func (s *server) Run() error {
	return http.ListenAndServe(":8080", s.router)
}

type Server interface {
	Run() error
	AddHandler(method, path string, handlerFn http.HandlerFunc)
}

var _ Server = &server{}
