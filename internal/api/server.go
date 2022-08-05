package api

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type Handler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type server struct {
	host   string
	router *chi.Mux
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func New(host string) *server {
	r := chi.NewRouter()

	r.Use(gzipHandle)

	return &server{
		router: r,
		host:   host,
	}
}

func (s *server) AddHandler(method, path string, handlerFn http.HandlerFunc) {
	s.router.MethodFunc(method, path, handlerFn)
}

func (s *server) Run() error {
	return http.ListenAndServe(s.host, s.router)
}

type Server interface {
	Run() error
	AddHandler(method, path string, handlerFn http.HandlerFunc)
}

var _ Server = &server{}
