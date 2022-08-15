package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/IgorAleksandroff/musthave-devops/internal/api"
	"github.com/IgorAleksandroff/musthave-devops/internal/pkg/metricscollection"
	"github.com/IgorAleksandroff/musthave-devops/utils/enviroment/serverconfig"
	"github.com/jackc/pgx/v4"
)

func main() {
	ctx, closeCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer closeCtx()

	config := serverconfig.Read()

	conn, err := pgx.Connect(ctx, config.AddressDB)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	metricsRepo := metricscollection.NewRepository(ctx, metricscollection.Config{
		StorePath:     config.StorePath,
		StoreInterval: config.StoreInterval,
		Restore:       config.Restore,
		AddressDB:     config.AddressDB,
	},
		conn,
	)
	metricsUC := metricscollection.NewUsecase(metricsRepo)

	server := api.New(config.Host, config.HashKey, metricsUC)

	metricsRepo.MemSync()

	log.Fatal(server.Run())
}
