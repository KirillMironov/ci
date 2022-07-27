package transport

import (
	"errors"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Handler is a handler for the HTTP requests.
type Handler struct {
	scheduler   scheduler
	logsStorage logsStorage
}

type (
	scheduler interface {
		Put(domain.Repository)
		Delete(domain.RepositoryURL)
	}
	logsStorage interface {
		GetById(id int) (domain.Log, error)
	}
)

// NewHandler creates a new Handler.
func NewHandler(scheduler scheduler, logsStorage logsStorage) *Handler {
	return &Handler{
		scheduler:   scheduler,
		logsStorage: logsStorage,
	}
}

// InitRoutes initializes the routes.
func (h Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	api := router.Group("/api/v1")
	{
		repository := api.Group("/repository")
		{
			repository.PUT("", h.putRepository)
			repository.DELETE("", h.deleteRepository)
		}
		log := api.Group("/log")
		{
			log.GET("/:id", h.getLogById)
		}
	}
	return router
}

// putRepository puts a new repository to the scheduler.
func (h Handler) putRepository(c *gin.Context) {
	var form struct {
		URL             string `json:"url" binding:"required"`
		Branch          string `json:"branch" binding:"required"`
		PollingInterval string `json:"polling_interval" binding:"required"`
	}

	err := c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	pollingInterval, err := time.ParseDuration(form.PollingInterval)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	h.scheduler.Put(domain.Repository{
		URL:             form.URL,
		Branch:          form.Branch,
		PollingInterval: pollingInterval,
	})
}

// deleteRepository deletes a repository from the scheduler.
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
