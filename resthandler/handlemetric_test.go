package resthandler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/internal/metricscollectionentity"
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

func Example_second() {
	metricsUC, _ := metricscollection.NewMetricsCollection(context.Background(), metricscollection.Config{})
	h := New(metricsUC, "hashKey")

	r := chi.NewRouter()
	r.HandleFunc("/update/", h.HandleJSONPost)
	r.HandleFunc("/value/", h.HandleJSONGet)

	metricSave := metricscollectionentity.Metrics{
		ID:    "name01",
		MType: "gauge",
		Value: func() *float64 { v := 0.1; return &v }(),
	}
	metricSave.CalcHash(fmt.Sprintf("%s:gauge:%f", metricSave.ID, *metricSave.Value), "hashKey")

	payloadSave, err := json.Marshal(metricSave)
	if err != nil {
		log.Println("payload marshal error")

		return
	}

	resSave := httptest.NewRecorder()
	reqSave, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewReader(payloadSave))
	reqSave.Header.Set(`Content-Type`, `application/json`)
	r.ServeHTTP(resSave, reqSave)
	resultSave := resSave.Result()
	fmt.Println(resultSave.StatusCode)
	resultSave.Body.Close()

	metricGet := metricscollectionentity.Metrics{
		ID:    "name01",
		MType: "gauge",
	}

	payloadGet, err := json.Marshal(metricGet)
	if err != nil {
		log.Println("payload marshal error")

		return
	}

	resGet := httptest.NewRecorder()
	reqGet, err := http.NewRequest(http.MethodPost, "/value/", bytes.NewReader(payloadGet))
	if err != nil {
		log.Println("create request error")
	}
	reqGet.Header.Set(`Content-Type`, `application/json`)
	r.ServeHTTP(resGet, reqGet)
	resultGet := resGet.Result()
	fmt.Println(resultGet.StatusCode)

	metric := metricscollectionentity.Metrics{}
	reader := json.NewDecoder(resultGet.Body)
	defer resultGet.Body.Close()
	reader.Decode(&metric)
	fmt.Println(*metric.Value)

	// Output:
	// 200
	// 200
	// 0.1
}
