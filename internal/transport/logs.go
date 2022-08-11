package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (h Handler) getLogById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	log, err := h.logsService.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, string(log.Data))
}
