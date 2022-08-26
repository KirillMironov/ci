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
	repositoriesUsecase domain.RepositoriesUsecase
	buildsUsecase       domain.BuildsUsecase
	logsUsecase         domain.LogsUsecase
}

func NewHandler(staticRootDir string, repositoriesUsecase domain.RepositoriesUsecase,
	buildsUsecase domain.BuildsUsecase, logsUsecase domain.LogsUsecase) *Handler {
	return &Handler{
		staticRootDir:       staticRootDir,
		repositoriesUsecase: repositoriesUsecase,
		buildsUsecase:       buildsUsecase,
		logsUsecase:         logsUsecase,
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
			repositories.DELETE("", h.deleteRepository)
			repositories.GET("", h.getRepositories)
			repositories.GET("/:repoId", h.getRepositoryById)
		}
		builds := api.Group("/repositories/:repoId/builds")
		{
			builds.GET("", h.getBuildsByRepositoryId)
			builds.GET("/:buildId", h.getBuild)
		}
		logs := api.Group("/logs")
		{
			logs.GET("/:buildId", h.getLogById)
		}
	}

	return router
}
