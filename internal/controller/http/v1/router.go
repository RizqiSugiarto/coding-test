package v1

import (
	"net/http"

	// Import generated swagger docs.
	_ "github.com/RizqiSugiarto/coding-test/docs"
	"github.com/RizqiSugiarto/coding-test/internal/usecase"
	"github.com/RizqiSugiarto/coding-test/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Cms API API
// @version     0.0.1
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(
	handler *gin.Engine,
	log logger.Interface,
	authUc usecase.Auth,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("api/v1")
	{
		newAuthRoutes(h, authUc, log)
	}
}
