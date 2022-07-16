package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type server struct {
	router *chi.Mux
}

func New() *server {
	r := chi.NewRouter()

	return &server{router: r}
}

func (s *server) AddHandler(method, path string, h Handler) {
	s.router.MethodFunc(method, path, h.Handle)
}

func (s *server) Run() error {
	return http.ListenAndServe(":8080", s.router)
}

type Server interface {
	Run() error
	AddHandler(method, path string, h Handler)
}

var _ Server = &server{}
