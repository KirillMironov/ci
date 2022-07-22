package transport

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Handler is a handler for the HTTP requests.
type Handler struct {
	scheduler scheduler
}

type scheduler interface {
	Put(domain.Repository)
	Delete(domain.RepositoryURL)
}

// NewHandler creates a new Handler.
func NewHandler(scheduler scheduler) *Handler {
	return &Handler{scheduler: scheduler}
}

// InitRoutes initializes the routes.
func (h Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	api := router.Group("/api/v1")
	{
		api.PUT("/repository", h.putRepository)
		api.DELETE("/repository", h.deleteRepository)
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
