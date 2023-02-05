package api

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/datacrypt"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/metrichandler"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server interface {
	Run()
}

type server struct {
	ctx        context.Context
	serverHTTP *http.Server
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func NewServer(ctx context.Context, cfg enviroment.ServerConfig, metricsUC metricscollection.MetricsCollection) *server {
	r := chi.NewRouter()

	r.Use(gzipUnzip)
	r.Use(gzipHandle)

	if cfg.CryptoKeyPath != "" {
		dc, err := datacrypt.New(
			datacrypt.WithPrivateKey(cfg.CryptoKeyPath),
			datacrypt.WithLabel("metrics"),
		)
		if err == nil && dc != nil {
			r.Use(dc.GetDecryptMiddleware())
		}
	}

	metricHandler := metrichandler.New(metricsUC, cfg.HashKey)

	r.MethodFunc(http.MethodPost, "/update/{TYPE}/{NAME}/{VALUE}", metricHandler.HandleMetricPost)
	r.MethodFunc(http.MethodGet, "/value/{TYPE}/{NAME}", metricHandler.HandleMetricGet)
	r.MethodFunc(http.MethodGet, "/", metricHandler.HandleMetricsGet)
	r.MethodFunc(http.MethodPost, "/update/", metricHandler.HandleJSONPost)
	r.MethodFunc(http.MethodPost, "/value/", metricHandler.HandleJSONGet)
	r.MethodFunc(http.MethodGet, "/ping", metricHandler.HandleDBPing)
	r.MethodFunc(http.MethodPost, "/updates/", metricHandler.HandleJSONPostBatch)

	return &server{
		ctx: ctx,
		serverHTTP: &http.Server{
			Handler:      r,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         cfg.Host,
		},
	}
}

func (s *server) Run() {
	notify := make(chan error, 1)
	go func() {
		notify <- s.serverHTTP.ListenAndServe()
		close(notify)
	}()

	select {
	case <-s.ctx.Done():
		log.Println("server interrupted by", s.ctx.Err())
	case err := <-notify:
		log.Printf("http server stopped: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	err := s.serverHTTP.Shutdown(ctx)
	if err != nil {
		log.Printf("error shutdown http server: %s", err)
	}
}

var _ Server = &server{}

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
