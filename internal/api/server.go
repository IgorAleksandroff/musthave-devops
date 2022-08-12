package api

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/IgorAleksandroff/musthave-devops/internal/api/metrichandler"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/go-chi/chi"
)

type Server interface {
	Run() error
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
func gzipUnzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if r.Header.Get(`Content-Encoding`) != `gzip` {
			// если не сжато методом gzip, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём *gzip.Reader, который будет читать тело запроса
		// и распаковывать его
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// не забывайте потом закрыть *gzip.Reader
		defer gz.Close()

		// меняем Body на тип gzip.Reader для распаковки данных
		r.Body = gz

		next.ServeHTTP(w, r)
	})
}

func New(host, key string, metricsUC metricscollection.Usecase) *server {
	r := chi.NewRouter()

	r.Use(gzipUnzip)
	r.Use(gzipHandle)

	metricHandler := metrichandler.New(metricsUC, key)

	r.MethodFunc(http.MethodPost, "/update/{TYPE}/{NAME}/{VALUE}", metricHandler.HandleMetricPost)
	r.MethodFunc(http.MethodGet, "/value/{TYPE}/{NAME}", metricHandler.HandleMetricGet)
	r.MethodFunc(http.MethodGet, "/", metricHandler.HandleMetricsGet)
	r.MethodFunc(http.MethodPost, "/update/", metricHandler.HandleJSONPost)
	r.MethodFunc(http.MethodPost, "/value/", metricHandler.HandleJSONGet)

	return &server{
		router: r,
		host:   host,
	}
}

func (s *server) Run() error {
	return http.ListenAndServe(s.host, s.router)
}

var _ Server = &server{}
