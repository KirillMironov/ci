package transport

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/gin-gonic/gin"
)

// Handler used to handle HTTP requests.
type Handler struct {
	scheduler           scheduler
	repositoriesService repositoriesService
	logsService         logsService
}

type (
	scheduler interface {
		Put(domain.Repository)
		Delete(id string)
	}
	repositoriesService interface {
		GetAll() ([]domain.Repository, error)
		GetById(id string) (domain.Repository, error)
	}
	logsService interface {
		GetById(id int) (domain.Log, error)
	}
)

func NewHandler(scheduler scheduler, repositoriesService repositoriesService, logsService logsService) *Handler {
	return &Handler{
		scheduler:           scheduler,
		repositoriesService: repositoriesService,
		logsService:         logsService,
	}
}

func (h Handler) Routes() *gin.Engine {
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
