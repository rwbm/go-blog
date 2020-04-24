package model

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	// ErrNoResults is used to indicate no results from a DB query
	ErrNoResults = errors.New("No results found")

	// ErrBadRequest (400) is returned for bad request (validation)
	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest)

	// ErrUnauthorized (401) is returned when user is not authorized
	ErrUnauthorized = echo.ErrUnauthorized
)
