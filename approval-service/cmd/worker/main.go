package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lht102/workflow-playground/approval-service/config"
	"github.com/lht102/workflow-playground/approval-service/entutil"
	"github.com/lht102/workflow-playground/approval-service/payment"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to init logger: %v\n", err)
	}

	zap.ReplaceGlobals(logger)

	temporalClient, err := client.Dial(config.GetTemporalClientOptions(cfg.TemporalConfig, logger))
	if err != nil {
		log.Panicf("Failed to create temporal client: %v\n", err)
	}
	defer temporalClient.Close()

	entClient, err := entutil.Open(cfg.MySQLConfig)
	if err != nil {
		log.Panicf("Failed to connect database: %v\n", err)
	}
	defer entClient.Close()

	workflowWorker := payment.NewWorkflowWorker(entClient, temporalClient, worker.Options{})
	if err := workflowWorker.Run(worker.InterruptCh()); err != nil {
		log.Panicf("Failed to start worker: %v\n", err)
	}
}
