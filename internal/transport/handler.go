package transport

import (
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/echox"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

// Handler used to handle HTTP requests.
type Handler struct {
	staticRootPath      string
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

func NewHandler(staticRootPath string, scheduler scheduler, repositoriesService repositoriesService,
	logsService logsService) *Handler {
	return &Handler{
		staticRootPath:      staticRootPath,
		scheduler:           scheduler,
		repositoriesService: repositoriesService,
		logsService:         logsService,
	}
}

func (h Handler) Routes() *echo.Echo {
	router := echo.New()
	router.Binder = echox.Binder{}
	router.Validator = echox.NewValidator()

	router.Use(
		middleware.Recover(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.DELETE, echo.OPTIONS},
		}),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  h.staticRootPath,
			HTML5: true,
			Skipper: func(c echo.Context) bool {
				return strings.HasPrefix(c.Request().URL.Path, "/api/")
			},
		}),
	)

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
