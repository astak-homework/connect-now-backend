package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	initDB()

	return &App{}
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
