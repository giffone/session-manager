package api

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"net/http"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/service"

	"github.com/labstack/echo/v4"
)

var logErr = "logErr"

type Handlers interface {
	CreateSessionOnCampus(c echo.Context) error
	CreateSessionOnPlatform(c echo.Context) error
	GetOnlineSessions(c echo.Context) error
	GetUserActivity(c echo.Context) error
}

type handlers struct {
	svc  service.Service
	rLog bool
}

func NewHandlers(logg echo.Logger, svc service.Service) Handlers {
	rLogStr := strings.ToLower(os.Getenv("REQ_LOG"))
	rLog := false
	if rLogStr == "true" {
		log.Println("request logging is enabled")
		rLog = true
	}
	return &handlers{
		svc:  svc,
		rLog: rLog,
	}
}

func (h *handlers) CreateSessionOnCampus(c echo.Context) error {
	var req request.Campus

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSessionOnCampus: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// validate data
	dto, err := req.Validate()
	if err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSessionOnCampus: validate: %s\n%s", err, req.Print()))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// create session on campus
	if sess, err := h.svc.CreateSessionOnCampus(c.Request().Context(), dto); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSessionOnCampus: %s\n%s", err, req.Print()))
		return customErrResponse(c, err, sess)
	}

	if h.rLog {
		log.Println(req.Print())
	}

	return created(c)
}

func (h *handlers) CreateSessionOnPlatform(c echo.Context) error {
	var req request.Platform

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSessionOnPlatform: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// validate data
	dto, err := req.Validate()
	if err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSessionOnPlatform: validate: %s\n%s", err, req.Print()))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// create session on platform
	if err := h.svc.CreateSessionOnPlatform(c.Request().Context(), dto); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSessionOnPlatform: %s\n%s", err, req.Print()))
		return customErrResponse(c, err, nil)
	}

	if h.rLog {
		log.Println(req.Print())
	}

	return created(c)
}

func (h *handlers) GetOnlineSessions(c echo.Context) error {
	sessions, err := h.svc.GetOnlineDashboard(c.Request().Context())
	if err != nil {
		c.Set(logErr, fmt.Sprintf("GetOnlineSessions: %s", err))
		return customErrResponse(c, err, nil)
	}

	c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")

	return ok(c, sessions)
}

func (h *handlers) GetUserActivity(c echo.Context) (err error) {
	var req request.UserActivity

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("GetUserActivity: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// validate data
	dto, err := req.Validate()
	if err != nil {
		c.Set(logErr, fmt.Sprintf("GetUserActivity: validate: %s\n%s", err, req.Print()))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	activity, err := h.svc.GetUserActivity(c.Request().Context(), dto)
	if err != nil {
		c.Set(logErr, fmt.Sprintf("GetUserActivity: %s\n%s", err, req.Print()))
		return customErrResponse(c, err, nil)
	}

	if h.rLog {
		log.Println(req.Print())
	}

	return ok(c, activity)
}

func customErrResponse(c echo.Context, err error, data any) error {
	defer printLogErr(c)

	if data == nil {
		data = []string{} // to show empty array
	}
	if errors.Is(err, response.ErrAccessDenied) {
		return c.JSON(http.StatusUnauthorized, response.Data{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    data,
		})
	}
	var errBadReq *response.ErrBadReq
	if errors.As(err, &errBadReq) {
		return c.JSON(http.StatusBadRequest, response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    data,
		})
	}

	return c.JSON(http.StatusInternalServerError, response.Data{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Data:    data,
	})
}

func printLogErr(c echo.Context) {
	if eLog := c.Get(logErr); eLog != nil {
		log.Printf("\n//----\n[error]: %v\n----\\\\\n", eLog)
	}
}

func created(c echo.Context) error {
	return c.JSON(http.StatusCreated, response.Data{
		Code:    http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    []string{},
	})
}

func ok(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, response.Data{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}
