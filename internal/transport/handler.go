package transport

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Handler is a handler for the HTTP requests.
type Handler struct {
	add chan<- domain.Repository
}

// NewHandler creates a new Handler.
func NewHandler(add chan<- domain.Repository) *Handler {
	return &Handler{add: add}
}

// InitRoutes initializes the routes.
func (h Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	api := router.Group("/api/v1")
	{
		api.PUT("/repository", h.addRepository)
	}
	return router
}

// addRepository adds a new repository to the scheduler.
func (h Handler) addRepository(c *gin.Context) {
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

	h.add <- domain.Repository{
		URL:             form.URL,
		Branch:          form.Branch,
		PollingInterval: pollingInterval,
	}
}
