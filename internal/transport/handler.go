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
	staticRootDir       string
	scheduler           scheduler
	repositoriesStorage domain.RepositoriesStorage
	buildsStorage       domain.BuildsStorage
	logsStorage         domain.LogsStorage
}

type scheduler interface {
	Add(domain.Repository)
	Remove(id string)
}

func NewHandler(staticRootDir string, s scheduler, rs domain.RepositoriesStorage, bs domain.BuildsStorage,
	ls domain.LogsStorage) *Handler {
	return &Handler{
		staticRootDir:       staticRootDir,
		scheduler:           s,
		repositoriesStorage: rs,
		buildsStorage:       bs,
		logsStorage:         ls,
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
			AllowMethods: []string{echo.GET, echo.POST, echo.DELETE, echo.OPTIONS},
		}),
		middleware.StaticWithConfig(middleware.StaticConfig{
			Root:  h.staticRootDir,
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
			repositories.POST("", h.addRepository)
			repositories.DELETE("", h.removeRepository)
			repositories.GET("", h.getRepositories)
			repositories.GET("/:repoId", h.getRepositoryById)
		}
		builds := api.Group("/repositories/:repoId/builds")
		{
			builds.GET("", h.getBuildsByRepoId)
			builds.GET("/:buildId", h.getBuildById)
		}
		logs := api.Group("/logs")
		{
			logs.GET("/:buildId", h.getLogById)
		}
	}

	return router
}
