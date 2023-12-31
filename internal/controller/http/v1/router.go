package v1

import (
	"os"

	"github.com/ABDURAZZAKK/avito_experiment/internal/service"

	log "github.com/sirupsen/logrus"

	"github.com/ABDURAZZAKK/avito_experiment/pkg/broker"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	STATIC_CSV_PATH = "assets/csv"
)

func NewRouter(handler *echo.Echo, services *service.Services, rabbit *broker.RabbitMQ) {
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())
	handler.Static("/assets/csv", STATIC_CSV_PATH)
	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })

	v1 := handler.Group("/api/v1")
	{
		newUserRoutes(v1.Group("/users"), services.User, rabbit)
		newSegmentRoutes(v1.Group("/segments"), services.Segment, services.User, rabbit)
		newFileRoutes(v1.Group("/stats"), rabbit)
	}
}

func setLogsFile() *os.File {
	file, err := os.OpenFile("./logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
