package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"

	"github.com/jtonynet/cine-catalogo/config"
	"github.com/jtonynet/cine-catalogo/internal/decorators"
	"github.com/jtonynet/cine-catalogo/internal/interfaces"
	"github.com/jtonynet/cine-catalogo/internal/logger"
	"github.com/jtonynet/cine-catalogo/internal/models"
)

type Database struct {
	DB  *gorm.DB
	log interfaces.Logger
}

func NewDatabase(cfg *config.Database) (*Database, error) {
	key := "database-init"

	l, err := logger.NewLogger()
	if err != nil {
		fmt.Printf("log database failure %v", err)
	}
	log := decorators.NewLoggerWithMetrics(l)

	log.WithField("origin", key).
		Info("database: trying open connection")

	strConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%v sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Pass,
		cfg.DB,
		cfg.Port)

	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(strConn))
	if err != nil {
		log.WithError(err).Error("database: error on connection")
		return nil, err
	}

	if cfg.MetricEnabled {
		db.Use(prometheus.New(prometheus.Config{
			DBName:          cfg.MetricDBName,          // `DBName` as metrics label
			RefreshInterval: cfg.MetricRefreshInterval, // refresh metrics interval (default 15 seconds)
			StartServer:     cfg.MetricStartServer,     // start http server to expose metrics
			HTTPServerPort:  cfg.MetricServerPort,      // configure http server port, default port 8080 (if you have configured multiple instances, only the first `HTTPServerPort` will be used to start server)
			MetricsCollector: []prometheus.MetricsCollector{
				&prometheus.Postgres{VariableNames: []string{"Threads_running"}},
			},
		}))
	}

	log.WithField("origin", key).
		Info("database: connection is openned")

	db.AutoMigrate(&models.Address{})
	db.AutoMigrate(&models.Cinema{})

	db.AutoMigrate(&models.Movie{})
	db.AutoMigrate(&models.Poster{})

	log.WithField("origin", key).
		Info("database: tables created")

	resultDatabase := &Database{DB: db, log: log}
	return resultDatabase, nil
}

func (db *Database) IsConnected() error {
	key := "database-is-connected"

	if err := db.DB.Raw("SELECT 1").Error; err != nil {
		db.log.WithError(err).
			WithField("origin", key).
			Error("error trying check readiness")
		return err
	}
	return nil
}
