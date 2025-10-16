package v1

import (
	"net/http"

	// Import generated swagger docs.
	_ "github.com/RizqiSugiarto/coding-test/docs"
	"github.com/RizqiSugiarto/coding-test/internal/controller/http/v1/middleware"
	"github.com/RizqiSugiarto/coding-test/internal/usecase"
	"github.com/RizqiSugiarto/coding-test/pkg/jwt"
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
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func NewRouter(
	handler *gin.Engine,
	log logger.Interface,
	authUc usecase.Auth,
	categoryUc usecase.Category,
	newsUc usecase.News,
	customPageUc usecase.CustomPage,
	commentUc usecase.Comment,
	jwtManager jwt.Manager,
) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Middleware
	authMiddleware := middleware.AuthMiddleware(jwtManager)

	// Routers
	h := handler.Group("api/v1")
	{
		newAuthRoutes(h, authUc, log)
		newCategoryRoutes(h, categoryUc, log, authMiddleware)
		newNewsRoutes(h, newsUc, log, authMiddleware)
		newCustomPageRoutes(h, customPageUc, log, authMiddleware)
		newCommentRoutes(h, commentUc, log)
	}
}
