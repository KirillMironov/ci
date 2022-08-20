package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/duration"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) addRepository(c echo.Context) error {
	var form struct {
		URL             string            `json:"url" validate:"required"`
		Branch          string            `json:"branch" validate:"required"`
		PollingInterval duration.Duration `json:"polling_interval" validate:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return err
	}

	err = h.repositoriesUsecase.Add(domain.Repository{
		URL:             form.URL,
		Branch:          form.Branch,
		PollingInterval: form.PollingInterval,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h Handler) deleteRepository(c echo.Context) error {
	var form struct {
		Id string `json:"id" validate:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return err
	}

	err = h.repositoriesUsecase.Delete(form.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h Handler) getRepositories(c echo.Context) error {
	repositories, err := h.repositoriesUsecase.GetAll()
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"repositories": repositories})
}

func (h Handler) getRepositoryById(c echo.Context) error {
	repository, err := h.repositoriesUsecase.GetById(c.Param("repoId"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, repository)
}
