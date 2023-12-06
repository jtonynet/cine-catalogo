package handlers

import (
	"github.com/jtonynet/cine-catalogo/internal/logger"
)

var (
	log *logger.Logger
)

func Init() {
	log = logger.NewLogger("handlers")
}
