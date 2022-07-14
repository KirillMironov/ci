package transport

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Handler is a handler for the HTTP requests.
type Handler struct {
	poller poller
}

// poller is a poller for the source code repository.
type poller interface {
	Start(domain.Repository)
}

// NewHandler creates a new Handler.
func NewHandler(poller poller) *Handler {
	return &Handler{poller: poller}
}

// InitRoutes initializes the routes.
func (h Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	api := router.Group("/api/v1")
	{
		api.POST("/repository", h.addRepository)
	}
	return router
}

// addRepository starts repository polling with a given interval.
func (h Handler) addRepository(c *gin.Context) {
	var form struct {
		URL             string `json:"url" binding:"required"`
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

	go h.poller.Start(domain.Repository{
		URL:             form.URL,
		PollingInterval: pollingInterval,
	})
}
