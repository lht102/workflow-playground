package entutil

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/lht102/workflow-playground/approval-service/config"
	"github.com/lht102/workflow-playground/approval-service/ent"
)

func Open(cfg config.MySQLConfig) (*ent.Client, error) {
	drv, err := sql.Open("mysql", config.GetMySQLDSN(cfg))
	if err != nil {
		return nil, fmt.Errorf("open mysql connection: %w", err)
	}

	db := drv.DB()
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Hour)

	return ent.NewClient(ent.Driver(drv)), nil
}
