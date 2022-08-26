package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) getBuildById(c echo.Context) error {
	build, err := h.buildsStorage.GetById(c.Param("buildId"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, build)
}

func (h Handler) getBuildsByRepoId(c echo.Context) error {
	builds, err := h.buildsStorage.GetAllByRepoId(c.Param("repoId"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"builds": builds})
}
