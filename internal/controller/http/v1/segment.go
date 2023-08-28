package v1

import (
	"net/http"

	"github.com/ABDURAZZAKK/avito_experiment/internal/service"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/broker"
	"github.com/labstack/echo/v4"
)

type segmentRoutes struct {
	segmentService service.Segment
	Rabbit         *broker.RabbitMQ
}

func newSegmentRoutes(g *echo.Group, segmentService service.Segment, rabbit *broker.RabbitMQ) *segmentRoutes {
	r := &segmentRoutes{
		segmentService: segmentService,
		Rabbit:         rabbit,
	}

	g.POST("/create", r.create)
	g.POST("/createAll", r.createAll)
	g.DELETE("/delete", r.delete)
	return r
}

type segmentCreateInput struct {
	Slug string `json:"slug"`
}

// @Summary Create segment
// @Description Create segment
// @Tags Segments
// @Accept json
// @Produce json
// @Success 201 {object} v1.segmentRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/segments/create [post]
func (r *segmentRoutes) create(c echo.Context) error {
	var input segmentCreateInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	slug, err := r.segmentService.Create(c.Request().Context(), input.Slug)
	if err != nil {
		if err == service.ErrAlreadyExists {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Slug string `json:"slug"`
	}

	return c.JSON(http.StatusCreated, response{
		Slug: slug,
	})
}

type segmentCreateAllInput struct {
	Slugs []string `json:"slugs"`
}

// @Summary Create segment
// @Description Create segment
// @Tags Segments
// @Accept json
// @Produce json
// @Success 201 {object} v1.segmentRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/segments/createAll [post]
func (r *segmentRoutes) createAll(c echo.Context) error {
	var input segmentCreateAllInput
	if err := c.Bind(&input); err != nil || len(input.Slugs) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	err := r.segmentService.CreateAll(c.Request().Context(), input.Slugs)
	if err != nil {
		if err == service.ErrAlreadyExists {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Message string `json:"message"`
	}

	return c.JSON(http.StatusCreated, response{
		Message: "Success",
	})
}

type deleteSegmentInput struct {
	Slug string `json:"slug"`
}

// @Summary Delete segment
// @Description Delete segment
// @Tags Segments
// @Accept json
// @Produce json
// @Success 200 {object} v1.segmentRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/segments/delete [delete]
func (r *segmentRoutes) delete(c echo.Context) error {
	var input deleteSegmentInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	slug, err := r.segmentService.Delete(c.Request().Context(), input.Slug)
	if err != nil {
		if err == service.ErrNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Slug string `json:"slug"`
	}
	return c.JSON(http.StatusOK, response{
		Slug: slug,
	})
}
