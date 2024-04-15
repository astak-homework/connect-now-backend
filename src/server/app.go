package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	authhttp "github.com/astak-homework/connect-now-backend/auth/delivery/http"
	authpostgres "github.com/astak-homework/connect-now-backend/auth/repository/postgresql"
	authusecase "github.com/astak-homework/connect-now-backend/auth/usecase"
	"github.com/astak-homework/connect-now-backend/config"
	"github.com/astak-homework/connect-now-backend/errors"
	"github.com/astak-homework/connect-now-backend/profile"
	profilehttp "github.com/astak-homework/connect-now-backend/profile/delivery/http"
	profilepostgres "github.com/astak-homework/connect-now-backend/profile/repository/postgresql"
	profileusecase "github.com/astak-homework/connect-now-backend/profile/usecase"
	"github.com/astak-homework/connect-now-backend/resources"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type App struct {
	httpServer *http.Server

	authUseCase    auth.UseCase
	profileUseCase profile.UseCase
}

func NewApp(cfg *config.Config) *App {
	conn := initDB(cfg.Postgres)

	loginRepo := authpostgres.NewLoginRepository(conn)
	profileRepo := profilepostgres.NewProfileRepository(conn)

	return &App{
		authUseCase:    authusecase.NewAuthUseCase(loginRepo, cfg.Auth),
		profileUseCase: profileusecase.NewProfileUseCase(profileRepo),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
		resources.Localize("./lang"),
		requestid.New(),
		errors.HttpResponseErrorHandler,
	)

	// Set up http handlers
	// SignUp/SignIn endpoints
	authhttp.RegisterHTTPEndpoints(router, a.authUseCase)

	// API endpoints
	user := router.Group("/user")

	profilehttp.RegisterHTTPEndpoints(user, a.authUseCase, a.profileUseCase)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB(cfg *config.Postgres) *pgxpool.Pool {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	dbConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to the database")
	}

	return conn
}
