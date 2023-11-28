package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gookit/validate"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
)

func NewHTTPError(code int, err any) HTTPError {
	return HTTPError{Code: code, Err: err}
}

type HTTPError struct {
	Code int
	Err  any
}

func (e HTTPError) Error() string {
	switch err := e.Err.(type) {
	case error:
		return err.Error()
	case fmt.Stringer:
		return err.String()
	case string:
		return err
	default:
		panic(fmt.Sprintf("can't get error message from %T", err))
	}
}

var ErrRequestMalformed = NewHTTPError(http.StatusBadRequest, "Request malformed")

func NewValidationErrorResponse(errs validate.Errors) HTTPError {
	return NewHTTPError(http.StatusBadRequest, map[string]any{"errors": errs})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Server struct {
	EchoInst      *echo.Echo
	PublicCtrl    PublicController
	ChallengeCtrl ChallengeController
}

func NewServer(i *do.Injector) Server {
	echoInst := echo.New()
	echoInst.Use(middleware.Logger())
	echoInst.Use(middleware.Recover())
	echoInst.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		httpErr, ok := err.(HTTPError)
		if !ok {
			httpErr = HTTPError{
				Code: http.StatusInternalServerError,
				Err:  http.StatusText(http.StatusInternalServerError),
			}
		}

		var msg any
		if e, ok := httpErr.Err.(json.Marshaler); ok {
			msg = e
		} else {
			msg = ErrorResponse{Error: httpErr.Error()}
		}

		// Send response
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(httpErr.Code)
		} else {
			err = c.JSON(httpErr.Code, msg)
		}
		if err != nil {
			echoInst.Logger.Error(err)
		}
	}

	return Server{
		EchoInst:      echoInst,
		PublicCtrl:    NewPublicController(echoInst.Group("")),
		ChallengeCtrl: NewChallengeController(echoInst.Group("/v1/challenges"), i),
	}
}

func (s Server) Serve(addr string) error {
	lst, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen to tcp port: %w", err)
	}
	log.Printf("Http server listening on %s", addr)
	return s.EchoInst.Server.Serve(lst)
}