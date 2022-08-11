package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/duration"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) putRepository(c echo.Context) error {
	var form struct {
		URL             string            `json:"url" binding:"required"`
		Branch          string            `json:"branch" binding:"required"`
		PollingInterval duration.Duration `json:"polling_interval" binding:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.scheduler.Put(domain.Repository{
		URL:             form.URL,
		Branch:          form.Branch,
		PollingInterval: form.PollingInterval,
	})

	return nil
}

func (h Handler) deleteRepository(c echo.Context) error {
	var form struct {
		Id string `json:"id" binding:"required"`
	}

	err := c.Bind(&form)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.scheduler.Delete(form.Id)

	return nil
}

func (h Handler) getRepositories(c echo.Context) error {
	type Repository struct {
		Id           string `json:"id"`
		URL          string `json:"url"`
		LatestCommit string `json:"latest_commit"`
	}

	var response struct {
		Repositories []Repository `json:"repositories"`
	}

	repositories, err := h.repositoriesService.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	for _, repo := range repositories {
		var r = Repository{
			Id:  repo.Id,
			URL: repo.URL,
		}
		if len(repo.Builds) > 0 {
			r.LatestCommit = repo.Builds[len(repo.Builds)-1].Commit.Hash
		}
		response.Repositories = append(response.Repositories, r)
	}

	return c.JSON(http.StatusOK, response)
}

func (h Handler) getRepositoryById(c echo.Context) error {
	repo, err := h.repositoriesService.GetById(c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, repo)
}
