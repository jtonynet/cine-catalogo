package handlers

import (
	"github.com/jtonynet/cine-catalogo/internal/decorators"
	"github.com/jtonynet/cine-catalogo/internal/interfaces"
	"github.com/jtonynet/cine-catalogo/internal/logger"
)

var (
	log interfaces.Logger
)

func Init() {
	l := logger.NewLogger("handlers")
	log = decorators.NewLoggerWithMetrics(l)
}
