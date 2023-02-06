package echo

import (
	"context"
	"fmt"
	nethttp "net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yazdsr/master-api/internal/pkg/logger"
	"github.com/yazdsr/master-api/internal/repository"
	"github.com/yazdsr/master-api/internal/transport/http"
	mdlware "github.com/yazdsr/master-api/internal/transport/http/middleware"
	"github.com/yazdsr/master-api/pkg/jwt"
)

type rest struct {
	echo             *echo.Echo
	adminController  *adminController
	userController   *userController
	adminMiddleware  *mdlware.AdminMiddleware
	serverController *serverController
}

func New(logger logger.Logger, repo repository.Postgres, jwt jwt.Jwt) http.Rest {
	return &rest{
		echo: echo.New(),
		adminController: &adminController{
			logger: logger,
			repo:   repo,
			jwt:    jwt,
		},
		adminMiddleware: &mdlware.AdminMiddleware{
			JwtPkg: jwt,
		},
		userController: &userController{
			logger: logger,
			repo:   repo,
		},
		serverController: &serverController{
			logger: logger,
			repo:   repo,
		},
	}
}

func (r *rest) Start(address string) error {
	r.echo.Use(middleware.Recover())
	r.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{nethttp.MethodGet, nethttp.MethodOptions, nethttp.MethodPut, nethttp.MethodPost, nethttp.MethodDelete, nethttp.MethodPatch},
		AllowCredentials: true,
	}))

	r.routing()
	return r.echo.Start(fmt.Sprintf(":%s", address))
}

func (r *rest) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.echo.Shutdown(ctx)
}
