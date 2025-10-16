package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RizqiSugiarto/coding-test/config"
	v1 "github.com/RizqiSugiarto/coding-test/internal/controller/http/v1"
	repoPg "github.com/RizqiSugiarto/coding-test/internal/repository/postgres"
	"github.com/RizqiSugiarto/coding-test/internal/usecase"
	"github.com/RizqiSugiarto/coding-test/pkg/httpserver"
	"github.com/RizqiSugiarto/coding-test/pkg/jwt"
	"github.com/RizqiSugiarto/coding-test/pkg/logger"
	pkgPg "github.com/RizqiSugiarto/coding-test/pkg/postgres"
	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.DBName,
	)

	pg, err := pkgPg.New(pgURL, pkgPg.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	jwtManager := jwt.NewJWTManager(&cfg.JWT)

	// Repo
	userRepo := repoPg.NewPostgresUserRepo(pg)
	categoryRepo := repoPg.NewPostgresCategoryRepo(pg)
	newsRepo := repoPg.NewPostgresNewsRepo(pg)
	customPageRepo := repoPg.NewPostgresCustomPageRepo(pg)

	// Usecase
	authUc := usecase.NewAuthUseCase(userRepo, jwtManager)
	categoryUc := usecase.NewCategoryUseCase(categoryRepo)
	newsUc := usecase.NewNewsUseCase(newsRepo)
	customPageUc := usecase.NewCustomPageUseCase(customPageRepo)

	initMigration(pgURL)

	if err = seedUsers(userRepo); err != nil {
		log.Error(fmt.Errorf("app - Run - seedUsers: %w", err))
	}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, log, authUc, categoryUc, newsUc, customPageUc, jwtManager)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))

		// Shutdown
		err = httpServer.Shutdown()
		if err != nil {
			log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
		}
	}
}
