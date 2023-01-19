package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"github.com/yazdsr/master-api/internal/config"
	"github.com/yazdsr/master-api/internal/pkg/logger/zap"
	"github.com/yazdsr/master-api/internal/repository/postgres"
	"github.com/yazdsr/master-api/internal/transport/http/echo"
	"github.com/yazdsr/master-api/pkg/jwt"
	"go.uber.org/zap/zapcore"
)

var serveCMD = &cli.Command{
	Name:    "serve",
	Aliases: []string{"s"},
	Usage:   "serve http",
	Action:  serve,
}

func serve(c *cli.Context) error {
	cfg := new(config.Config)
	config.ReadEnv(cfg)
	f, err := os.OpenFile("logs/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logger := zap.New(f, zapcore.ErrorLevel)

	// create repo instance
	repo, err := postgres.New(cfg.Psql, logger)

	if err != nil {
		return err
	}

	jwt := jwt.New(cfg.App.Secret)

	restServer := echo.New(logger, repo, jwt)
	go func() {
		if err := restServer.Start(cfg.App.Port); err != nil {
			logger.Error(fmt.Sprintf("error happen while serving: %v", err))
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	fmt.Println("\nReceived an interrupt, closing connections...")

	if err := restServer.Shutdown(); err != nil {
		fmt.Println("\nRest server doesn't shutdown in 10 seconds")
	}

	return nil
}
