package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/logger"
	"github.com/jtonynet/cine-catalogo/models"
)

var (
	DB  *gorm.DB
	err error
	l   *logger.Logger
)

func Init(cfg config.Database) error {
	l = logger.NewLogger("database")

	l.Info("database: trying open connection")

	strConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Pass,
		cfg.DB,
		cfg.Port)

	DB, err = gorm.Open(postgres.Open(strConn))
	if err != nil {
		l.Errorf("database: error on connection %v", err.Error())
		return err
	}

	l.Info("database: connection is openned")

	DB.AutoMigrate(&models.Address{})
	DB.AutoMigrate(&models.Movie{})

	l.Info("database: tables created")

	return nil
}

func IsConnected() error {
	if err := DB.Raw("SELECT 1").Error; err != nil {
		l.Errorf("database: error trying check readiness %v", err.Error())
		return err
	}
	return nil
}
