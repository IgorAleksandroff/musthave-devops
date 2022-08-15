package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorymemo"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection/repositorypg"
	"github.com/IgorAleksandroff/musthave-devops/utils/enviroment/serverconfig"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeCtx()

	config := serverconfig.Read()

	var conn *pgxpool.Pool
	var err error
	if config.AddressDB != "" {
		conn, err = pgxpool.Connect(ctx, config.AddressDB)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
			os.Exit(1)
		}
		log.Printf("connect to DB: %v", conn.Config())
		defer conn.Close()
	}

	metricsRepo := repositorymemo.NewRepository(ctx, repositorymemo.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
	})
	connectionTester := repositorypg.NewPinger(ctx, conn)
	metricsUC := metricscollection.NewUsecase(metricsRepo)

	server := api.New(config.Host, config.HashKey, metricsUC, connectionTester)

	metricsRepo.MemSync()

	log.Fatal(server.Run())
}
