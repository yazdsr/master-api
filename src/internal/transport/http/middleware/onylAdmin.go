package middleware

import (
	"strings"

	"github.com/yazdsr/master-api/internal/transport/http/response"

	"github.com/labstack/echo/v4"
	"github.com/yazdsr/master-api/pkg/jwt"
)

type AdminMiddleware struct {
	JwtPkg jwt.Jwt
}

func NewUserMiddleware(jwtPkg jwt.Jwt) AdminMiddleware {
	return AdminMiddleware{
		JwtPkg: jwtPkg,
	}
}

func (am *AdminMiddleware) OnlyAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		errResp := response.Error{
			Code:    echo.ErrUnauthorized.Code,
			Message: "UnAuthorized Access",
		}

		if auth == "" {
			return c.JSON(errResp.Code, errResp)
		}
		token := strings.Split(auth, " ")
		if len(token) != 2 {
			return c.JSON(errResp.Code, errResp)
		}

		if token[0] != "Bearer" {
			return c.JSON(errResp.Code, errResp)
		}

		claims, err := am.JwtPkg.ParseToken(token[1])
		if err != nil {
			return c.JSON(errResp.Code, errResp)
		}
		uID, ok := claims["sub"]
		if !ok {
			return c.JSON(errResp.Code, errResp)
		}
		userID := int(uID.(float64))
		c.Set("user_id", int(userID))
		return next(c)
	}
}
