package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lht102/workflow-playground/approval-service/config"
	"github.com/lht102/workflow-playground/approval-service/entutil"
	"github.com/lht102/workflow-playground/approval-service/http/rest"
	"github.com/lht102/workflow-playground/approval-service/payment"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to init logger: %v\n", err)
	}

	zap.ReplaceGlobals(logger)

	stopSignalCh := make(chan os.Signal, 1)
	signal.Notify(stopSignalCh, os.Interrupt, syscall.SIGTERM)

	temporalClient, err := client.Dial(config.GetTemporalClientOptions(cfg.TemporalConfig, logger))
	if err != nil {
		logger.Fatal("Failed to create temporal client", zap.Error(err))
	}
	defer temporalClient.Close()

	entClient, err := entutil.Open(cfg.MySQLConfig)
	if err != nil {
		logger.Fatal("Failed to connect database", zap.Error(err))
	}
	defer entClient.Close()

	paymentService := payment.NewService(entClient, temporalClient)
	httpServer := rest.NewServer(paymentService, cfg.RestPort, logger)

	go func() {
		logger.Sugar().Infof("Listening on port %v", cfg.RestPort)

		if err := httpServer.Open(); err != nil {
			logger.Fatal("Failed to listen and serve for http server", zap.Error(err))
		}
	}()

	<-stopSignalCh

	if err := httpServer.Close(); err != nil {
		logger.Fatal("Failed to shutdown http server", zap.Error(err))
	}

	logger.Info("Done")
}
