package echo

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/internal/pkg/logger"
	"github.com/yazdsr/master-api/internal/repository"
	"github.com/yazdsr/master-api/internal/transport/http/response"
)

type userController struct {
	logger logger.Logger
	repo   repository.Postgres
}

func (uc *userController) FindAllUsers(c echo.Context) error {
	users, err := uc.repo.FindAllUsers()
	if err != nil {
		rErr := response.Error{
			Code:    err.StatusCode(),
			Message: err.Error(),
			Errors:  err.Errors(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	return c.JSON(http.StatusOK, users)
}

func (uc *userController) FindUserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		}
		return c.JSON(rErr.Code, rErr)
	}
	user, rerr := uc.repo.FindUserByID(id)
	if err != nil {
		rErr := response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	return c.JSON(http.StatusOK, user)
}

func (uc *userController) CreateUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Request",
		}
		return c.JSON(rErr.Code, rErr)
	}
	if user.Username == "" || user.Password == "" || user.FullName == "" || user.ServerID == 0 || user.ValidUntil.IsZero() {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Fill Required Fields",
		}
		return c.JSON(rErr.Code, rErr)
	}

	// send json-rpc call to server to create user

	err := uc.repo.CreateUser(*user)
	if err != nil {
		rErr := response.Error{
			Code:    err.StatusCode(),
			Message: err.Error(),
			Errors:  err.Errors(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	return c.JSON(http.StatusCreated, response.Success{
		Code:    http.StatusCreated,
		Message: "User Created",
	})
}

func (uc *userController) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		}
		return c.JSON(rErr.Code, rErr)
	}
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Request",
		}
		return c.JSON(rErr.Code, rErr)
	}
	if id == 0 || user.Username == "" || user.Password == "" || user.FullName == "" || user.ServerID == 0 || user.ValidUntil.IsZero() {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Fill Required Fields",
		}
		return c.JSON(rErr.Code, rErr)
	}

	// send json-rpc call to server to update user

	user.ID = id
	rerr := uc.repo.UpdateUser(*user)
	if err != nil {
		rErr := response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	return c.JSON(http.StatusCreated, response.Success{
		Code:    http.StatusCreated,
		Message: "User Updated",
	})
}

func (uc *userController) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		}
		return c.JSON(rErr.Code, rErr)
	}

	// send json-rpc call to server to delete user

	rerr := uc.repo.DeleteUser(id)
	if err != nil {
		rErr := response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		}
		return c.JSON(rErr.Code, rErr)
	}
	return c.JSON(http.StatusOK, response.Success{
		Code:    http.StatusOK,
		Message: "User Deleted",
	})
}
