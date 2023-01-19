package echo

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yazdsr/master-api/internal/pkg/logger"
	"github.com/yazdsr/master-api/internal/repository"
	"github.com/yazdsr/master-api/internal/transport/http/request"
	"github.com/yazdsr/master-api/internal/transport/http/response"
	"github.com/yazdsr/master-api/pkg/jwt"
)

type adminController struct {
	logger logger.Logger
	repo   repository.Postgres
	jwt    jwt.Jwt
}

func (ac *adminController) login(c echo.Context) error {
	req := new(request.LoginRequest)
	if err := c.Bind(req); err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Request",
		}
		return c.JSON(rErr.Code, rErr)
	}
	if req.Username == "" || req.Password == "" {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Username or Password is empty",
		}
		return c.JSON(rErr.Code, rErr)
	}
	admin, rerr := ac.repo.FindAdminByUsernameAndPAssword(req.Username, req.Password)
	if rerr != nil {
		rErr := response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	if admin == nil || admin.ID == 0 {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Credentials",
		}
		return c.JSON(rErr.Code, rErr)
	}
	claims := map[string]interface{}{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat": time.Now().Unix(),
	}
	token, err := ac.jwt.GenerateToken(claims)
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	return c.JSON(http.StatusOK, response.Login{
		Token: token,
	})
}
