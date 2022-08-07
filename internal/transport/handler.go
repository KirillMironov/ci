package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/duration"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// Handler is a handler for the HTTP requests.
type Handler struct {
	scheduler           scheduler
	logsStorage         logsStorage
	repositoriesStorage repositoriesStorage
}

type (
	scheduler interface {
		Put(domain.Repository)
		Delete(domain.RepositoryURL)
	}
	logsStorage interface {
		GetById(id int) (domain.Log, error)
	}
	repositoriesStorage interface {
		GetAll() ([]domain.Repository, error)
		GetById(id string) (domain.Repository, error)
	}
)

// NewHandler creates a new Handler.
func NewHandler(scheduler scheduler, ls logsStorage, rs repositoriesStorage) *Handler {
	return &Handler{
		scheduler:           scheduler,
		logsStorage:         ls,
		repositoriesStorage: rs,
	}
}

// InitRoutes initializes the routes.
func (h Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), h.corsMiddleware)
	api := router.Group("/api/v1")
	{
		repositories := api.Group("/repositories")
		{
			repositories.PUT("", h.putRepository)
			repositories.DELETE("", h.deleteRepository)
			repositories.GET("", h.getRepositories)
			repositories.GET("/:id", h.getRepositoryById)
		}
		logs := api.Group("/logs")
		{
			logs.GET("/:id", h.getLogById)
		}
	}
	return router
}

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
		URL string `json:"url" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	h.scheduler.Delete(domain.RepositoryURL(form.URL))
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

	repositories, err := h.repositoriesStorage.GetAll()
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
	repo, err := h.repositoriesStorage.GetById(c.Param("id"))
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

func (h Handler) getLogById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	log, err := h.logsStorage.GetById(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, string(log.Data))
}
