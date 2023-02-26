package api

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"

	"github.com/IgorAleksandroff/musthave-devops/enviroment"
	"github.com/IgorAleksandroff/musthave-devops/internal/datacrypt"
	"github.com/IgorAleksandroff/musthave-devops/internal/generated/rpc"
	"github.com/IgorAleksandroff/musthave-devops/internal/grpchandler"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/resthandler"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 3 * time.Second
)

type Server interface {
	Run(ctx context.Context)
}

type server struct {
	serverHTTP   *http.Server
	serverGRPC   *grpc.Server
	gRPCListener net.Listener
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func NewServer(cfg enviroment.ServerConfig, metricsUC metricscollection.MetricsCollection) (*server, error) {
	// init HTTP server
	r := chi.NewRouter()

	r.Use(gzipUnzip)
	r.Use(gzipHandle)
	r.Use(cfg.GetTrustedIPMiddleware())

	if cfg.CryptoKeyPath != "" {
		dc, err := datacrypt.New(
			datacrypt.WithPrivateKey(cfg.CryptoKeyPath),
			datacrypt.WithLabel("metrics"),
		)
		if err == nil && dc != nil {
			r.Use(dc.GetDecryptMiddleware())
		}
	}

	metricRESTHandler := resthandler.New(metricsUC, cfg.HashKey)

	r.MethodFunc(http.MethodPost, "/update/{TYPE}/{NAME}/{VALUE}", metricRESTHandler.HandleMetricPost)
	r.MethodFunc(http.MethodGet, "/value/{TYPE}/{NAME}", metricRESTHandler.HandleMetricGet)
	r.MethodFunc(http.MethodGet, "/", metricRESTHandler.HandleMetricsGet)
	r.MethodFunc(http.MethodPost, "/update/", metricRESTHandler.HandleJSONPost)
	r.MethodFunc(http.MethodPost, "/value/", metricRESTHandler.HandleJSONGet)
	r.MethodFunc(http.MethodGet, "/ping", metricRESTHandler.HandleDBPing)
	r.MethodFunc(http.MethodPost, "/updates/", metricRESTHandler.HandleJSONPostBatch)

	// init GRPC server
	listen, err := net.Listen("tcp", cfg.GRPSSocket)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
		grpcValidator.UnaryServerInterceptor(),
		cfg.GetTrustedIPInterceptor(),
	)))
	metricGRPCHandler := grpchandler.New(metricsUC, cfg.HashKey)
	rpc.RegisterMetricsCollectionServer(s, metricGRPCHandler)

	return &server{
		serverHTTP: &http.Server{
			Handler:      r,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         cfg.Host,
		},
		serverGRPC:   s,
		gRPCListener: listen,
	}, nil
}

func (s server) Run(ctx context.Context) {
	notifyHTTP := make(chan error, 1)
	go func() {
		notifyHTTP <- s.serverHTTP.ListenAndServe()
		close(notifyHTTP)
	}()

	notifyGRPC := make(chan error, 1)
	go func() {
		notifyGRPC <- s.serverGRPC.Serve(s.gRPCListener)
		close(notifyGRPC)
	}()

	select {
	case <-ctx.Done():
		log.Println("server interrupted by", ctx.Err())

		s.serverGRPC.GracefulStop()
		s.shutdownHTTP()
	case err := <-notifyHTTP:
		log.Printf("HTTP server stopped: %s", err)

		s.serverGRPC.GracefulStop()
	case err := <-notifyGRPC:
		log.Printf("gRPC server stopped: %s", err)

		s.shutdownHTTP()
	}
}

var _ Server = &server{}

func (s server) shutdownHTTP() {
	ctxShutdown, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	err := s.serverHTTP.Shutdown(ctxShutdown)
	if err != nil {
		log.Printf("error shutdown http server: %s", err)
	}
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
