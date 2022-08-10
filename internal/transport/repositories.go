package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/duration"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) putRepository(c *gin.Context) {
	var form struct {
		URL             string            `json:"url" binding:"required"`
		Branch          string            `json:"branch" binding:"required"`
		PollingInterval duration.Duration `json:"polling_interval" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	h.scheduler.Put(domain.Repository{
		URL:             form.URL,
		Branch:          form.Branch,
		PollingInterval: form.PollingInterval,
	})
}

func (h Handler) deleteRepository(c *gin.Context) {
	var form struct {
		Id string `json:"id" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	h.scheduler.Delete(form.Id)
}

func (h Handler) getRepositories(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
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

	c.JSON(http.StatusOK, response)
}

func (h Handler) getRepositoryById(c *gin.Context) {
	repo, err := h.repositoriesService.GetById(c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, repo)
}
