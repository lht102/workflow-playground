package main

import (
	"context"
	"log"
	"os"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lht102/workflow-playground/approval-service/config"
	"github.com/lht102/workflow-playground/approval-service/entutil"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	dir, err := migrate.NewLocalDir("./migrations")
	if err != nil {
		log.Panicf("Failed to create atlas migration directory: %v\n", err)
	}

	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithFormatter(sqltool.GolangMigrateFormatter),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}

	if len(os.Args) != 2 {
		log.Fatalln("Migration name is required. Run the command with <name>.")
	}

	client, err := entutil.Open(cfg.MySQLConfig)
	if err != nil {
		log.Fatalf("Failed to connect database: %v\n", err)
	}
	defer client.Close()

	if err := client.Schema.NamedDiff(ctx, os.Args[1], opts...); err != nil {
		log.Panicf("Failed to create schema resources: %v\n", err)
	}
}
