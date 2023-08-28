package v1

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ABDURAZZAKK/avito_experiment/pkg/broker"
	"github.com/labstack/echo/v4"
)

type fileRoutes struct {
	Rabbit *broker.RabbitMQ
}

func newFileRoutes(g *echo.Group, rabbit *broker.RabbitMQ) {
	r := &fileRoutes{
		Rabbit: rabbit,
	}
	g.POST("/createCSVPerStats", r.createCSVFromUsersSegments)
}

type createCSVInput struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func (r *fileRoutes) createCSVFromUsersSegments(c echo.Context) error {
	var input createCSVInput
	if err := c.Bind(&input); err != nil || input.Month <= 0 || input.Month > 12 {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	filename := fmt.Sprintf("%s/user_segments_%v.csv",
		STATIC_CSV_PATH,
		time.Date(
			input.Year,
			time.Month(input.Month),
			1, 0, 0, 0, 0, time.Local).
			Format("2006_01_02"))

	msg, err := broker.MsgSerialize(broker.Message{
		"task":     "createCSVFromUsersSegments",
		"year":     input.Year,
		"month":    input.Month,
		"filename": filename,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	err = r.Rabbit.Publish(msg)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		URL string `json:"url"`
	}
	return c.JSON(http.StatusOK, response{
		URL: fmt.Sprintf("http://localhost:8000/%s", filename),
	})

}
