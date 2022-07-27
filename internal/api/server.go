package api

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

const EnvServerURL = "ADDRESS"
const DefaultServerURL = ":8080"

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
	return http.ListenAndServe(getEnvString(EnvServerURL, DefaultServerURL), s.router)
}

type Server interface {
	Run() error
	AddHandler(method, path string, h Handler)
}

var _ Server = &server{}

func getEnvString(envName, defaultValue string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Println("empty env")
		return defaultValue
	}
	return value
}
