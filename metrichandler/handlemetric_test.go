package metrichandler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
	"github.com/go-chi/chi"
)

func Example() {
	metricsUC, _ := metricscollection.NewMetricsCollection(context.Background(), metricscollection.Config{})
	h := New(metricsUC, "hashKey")

	urls := []string{
		"/update/gauge/name01/0.1",
		"/update/counter/name02/2",
	}

	r := chi.NewRouter()
	r.HandleFunc("/update/{TYPE}/{NAME}/{VALUE}", h.HandleMetricPost)
	r.HandleFunc("/value/{TYPE}/{NAME}/", h.HandleMetricGet)
	r.HandleFunc("/", h.HandleMetricsGet)

	for _, url := range urls {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, url, nil)
		r.ServeHTTP(res, req)
		result := res.Result()
		fmt.Println(result.StatusCode)
		result.Body.Close()
	}

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/value/gauge/name01/", nil)
	r.ServeHTTP(res, req)
	result := res.Result()
	value, _ := io.ReadAll(result.Body)
	defer result.Body.Close()
	fmt.Println(string(value))

	// меняем текущую директорию, чтобы был доступен html шаблон
	os.Chdir("../")
	resMetrics := httptest.NewRecorder()
	reqMetrics, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(resMetrics, reqMetrics)
	fmt.Println(resMetrics.Header().Get("Content-Type"))

	// Output:
	// 200
	// 200
	// 0.1
	// text/html
}
