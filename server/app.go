package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/astak-homework/connect-now-backend/auth"
	authpostgres "github.com/astak-homework/connect-now-backend/auth/repository/postgresql"
	authusecase "github.com/astak-homework/connect-now-backend/auth/usecase"
	"github.com/astak-homework/connect-now-backend/profile"
	profilepostgres "github.com/astak-homework/connect-now-backend/profile/repository/postgresql"
	profileusecase "github.com/astak-homework/connect-now-backend/profile/usecase"
	"github.com/gin-gonic/gin"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type App struct {
	httpServer *http.Server

	authUseCase    auth.UseCase
	profileUseCase profile.UseCase
}

func NewApp() *App {
	conn := initDB()

	loginRepo := authpostgres.NewLoginRepository(conn, viper.GetString("postgresql.login_table"))
	profileRepo := profilepostgres.NewProfileRepository(conn, viper.GetString("postgresql.profile_table"))

	return &App{
		authUseCase:    authusecase.NewAuthUseCase(loginRepo, viper.GetString("auth.hash_salt"), []byte(viper.GetString("auth.signing_key")), viper.GetDuration("auth.token_ttl")),
		profileUseCase: profileusecase.NewProfileUseCase(profileRepo),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

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

func initDB() *pgxpool.Pool {
	dbConfig, err := pgxpool.ParseConfig(viper.GetString("postgresql.uri"))
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
