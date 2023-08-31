package v1

import (
	"net/http"

	"github.com/ABDURAZZAKK/avito_experiment/internal/service"
	"github.com/ABDURAZZAKK/avito_experiment/pkg/broker"

	"github.com/labstack/echo/v4"
)

type userRoutes struct {
	userService service.User
	Rabbit      *broker.RabbitMQ
}

func newUserRoutes(g *echo.Group, userService service.User, rabbit *broker.RabbitMQ) {
	r := &userRoutes{
		userService: userService,
		Rabbit:      rabbit,
	}
	g.GET("", r.get)
	g.GET("/segments", r.getSegments)
	g.POST("/create", r.create)
	g.POST("/addSegments", r.addSegments)
	g.DELETE("/delete", r.delete)

}

type userCreateInput struct {
	Slug string `json:"slug"`
}

// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} v1.userRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/users/create [post]
func (r *userRoutes) create(c echo.Context) error {
	var input userCreateInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	id, err := r.userService.Create(c.Request().Context(), input.Slug)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: id,
	})
}

type getUserInput struct {
	Id int `query:"id"`
}

// @Summary Get user
// @Description Get user
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} v1.userRoutes.create.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/users [get]
func (r *userRoutes) get(c echo.Context) error {
	var input getUserInput
	if err := c.Bind(&input); err != nil || input.Id <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	user, err := r.userService.GetById(c.Request().Context(), input.Id)
	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Id   int    `json:"id"`
		Slug string `json:"slug"`
	}
	return c.JSON(http.StatusOK, response{
		Id:   user.Id,
		Slug: user.Slug,
	})
}

type getUserSegmentsInput struct {
	Id int `query:"id"`
}

func (r *userRoutes) getSegments(c echo.Context) error {
	var input getUserSegmentsInput
	if err := c.Bind(&input); err != nil || input.Id <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	segments, err := r.userService.GetSegments(c.Request().Context(), input.Id)
	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Segments []string `json:"segments"`
	}
	return c.JSON(http.StatusOK, response{
		Segments: segments,
	})
}

type changeUserSegmentsInput struct {
	Id         int      `json:"id"`
	AddList    []string `json:"add_list"`
	RemoveList []string `json:"remove_list"`
	DeleteAt   string   `json:"delete_at,omitempty"`
}

func (r *userRoutes) addSegments(c echo.Context) error {
	var input changeUserSegmentsInput
	if err := c.Bind(&input); err != nil ||
		input.Id <= 0 ||
		len(input.AddList) == 0 && len(input.RemoveList) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	err := r.userService.ChangeSegments(c.Request().Context(), input.Id, input.AddList, input.RemoveList)
	if err != nil {
		if err == service.ErrAlreadyExists {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		if err == service.ErrNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	if input.DeleteAt != "" {
		msg, err := broker.MsgSerialize(broker.Message{
			"task":     "DeleteSegmentFromUserOnTime",
			"time":     input.DeleteAt,
			"user":     input.Id,
			"segments": input.AddList,
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
	}

	type response struct {
		Message string `json:"message"`
	}

	return c.JSON(http.StatusCreated, response{
		Message: "Success",
	})
}

type deleteUserInput struct {
	Id int `json:"id"`
}

func (r *userRoutes) delete(c echo.Context) error {
	var input deleteUserInput
	if err := c.Bind(&input); err != nil || input.Id <= 0 {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}
	id, err := r.userService.Delete(c.Request().Context(), input.Id)
	if err != nil {
		if err == service.ErrUserNotFound {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return err
		}
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}
	type response struct {
		Id int `json:"id"`
	}
	return c.JSON(http.StatusOK, response{
		Id: id,
	})
}
