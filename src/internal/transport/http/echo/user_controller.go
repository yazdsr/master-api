package echo

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/labstack/echo/v4"
	"github.com/yazdsr/master-api/internal/pkg/logger"
	"github.com/yazdsr/master-api/internal/repository"
	"github.com/yazdsr/master-api/internal/transport/http/request"
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
	user := new(request.CreateUser)
	if err := c.Bind(user); err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Request",
		}
		return c.JSON(rErr.Code, rErr)
	}
	if user.Username == "" || user.Password == "" || user.FullName == "" || user.ServerID == 0 || user.ValidUntil.IsZero() || user.StartDate.IsZero() {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Fill Required Fields",
		}
		return c.JSON(rErr.Code, rErr)
	}
	if user.StartDate.Unix() >= user.ValidUntil.Unix() {
		return c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Message: "Start Date should be greater than due date",
		})
	}

	server, rErr := uc.repo.FindServerByID(user.ServerID)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	_, rErr = uc.repo.FindUserByUsernameAndServerID(user.Username, user.ServerID)
	if rErr == nil {
		return c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Message: "User Already Exists",
		})
	}
	if rErr.StatusCode() == http.StatusInternalServerError {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	users, rerr := uc.repo.FindAllUsers()
	if rerr != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	for _, u := range users {
		if u.ServerID == user.ServerID {
			if (user.StartDate.Unix() < u.ValidUntil.Unix() && user.StartDate.Unix() > u.StartDate.Unix()) || (user.ValidUntil.Unix() > u.StartDate.Unix() && user.ValidUntil.Unix() < u.ValidUntil.Unix()) {
				return c.JSON(http.StatusBadRequest, response.Error{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("There is start and due date conflict between this user and user with id %d and username %s", u.ID, u.Username),
				})
			}
		}
	}

	url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
	params := &request.CreateUserRPC{
		Username: user.Username,
		Password: user.Password,
	}

	msg, err := json2.EncodeClientRequest("createUser", params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error in encoding request",
		})
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error in sending request to %s: %s", url, err.Error()),
		})
	}
	defer resp.Body.Close()

	var rpcErr json2.Error
	err = json2.DecodeClientResponse(resp.Body, &rpcErr)
	jErr, ok := err.(*json2.Error)
	if ok {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: jErr.Error(),
		})
	}

	if rErr = uc.repo.CreateUser(*user); rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
			Errors:  rErr.Errors(),
		})
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
	user := new(request.UpdateUser)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid Request",
		})
	}
	if id == 0 || user.FullName == "" || user.ValidUntil.IsZero() || user.StartDate.IsZero() {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Fill Required Fields",
		}
		return c.JSON(rErr.Code, rErr)
	}

	usr, rerr := uc.repo.FindUserByID(id)
	if rerr != nil {
		return c.JSON(rerr.StatusCode(), response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
		})
	}

	server, rerr := uc.repo.FindServerByID(usr.ServerID)
	if rerr != nil {
		return c.JSON(rerr.StatusCode(), response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
		})
	}

	users, rerr := uc.repo.FindAllUsers()
	if rerr != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	for _, u := range users {
		if u.ServerID == usr.ServerID && u.ID != usr.ID {
			if (user.StartDate.Unix() < u.ValidUntil.Unix() && user.StartDate.Unix() > u.StartDate.Unix()) || (user.ValidUntil.Unix() > u.StartDate.Unix() && user.ValidUntil.Unix() < u.ValidUntil.Unix()) {
				return c.JSON(http.StatusBadRequest, response.Error{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("There is start and due date conflict between this user and user with id %d and username %s", u.ID, u.Username),
				})
			}
		}
	}

	if user.Password != "" {
		url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
		params := &request.UpdateUserRPC{
			Username: usr.Username,
			Password: user.Password,
		}

		msg, err := json2.EncodeClientRequest("updatePassword", params)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: "Error in encoding request",
			})
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			})
		}

		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Error in sending request to %s: %s", url, err.Error()),
			})
		}
		defer resp.Body.Close()

		var rpcErr json2.Error
		err = json2.DecodeClientResponse(resp.Body, &rpcErr)
		jErr, ok := err.(*json2.Error)
		if ok {
			return c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: jErr.Error(),
			})
		}
	}

	rerr = uc.repo.UpdateUser(id, *user)
	if rerr != nil {
		return c.JSON(rerr.StatusCode(), response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		})
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

	user, rErr := uc.repo.FindUserByID(id)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	server, rErr := uc.repo.FindServerByID(user.ServerID)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
	params := &request.DeleteUserRPC{
		Username: user.Username,
	}

	msg, err := json2.EncodeClientRequest("deleteUser", params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error in encoding request",
		})
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error in sending request to %s: %s", url, err.Error()),
		})
	}
	defer resp.Body.Close()

	var rpcErr json2.Error
	err = json2.DecodeClientResponse(resp.Body, &rpcErr)
	jErr, ok := err.(*json2.Error)
	if ok {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: jErr.Error(),
		})
	}

	rerr := uc.repo.DeleteUser(id)
	if rerr != nil {
		return c.JSON(rerr.StatusCode(), response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		})
	}
	return c.JSON(http.StatusOK, response.Success{
		Code:    http.StatusOK,
		Message: "User Deleted",
	})
}

func (uc *userController) ActiveUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		}
		return c.JSON(rErr.Code, rErr)
	}

	user, rErr := uc.repo.FindUserByID(id)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	server, rErr := uc.repo.FindServerByID(user.ServerID)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
	params := &request.ActiveUserRPC{
		Username: user.Username,
	}

	msg, err := json2.EncodeClientRequest("activeUser", params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error in encoding request",
		})
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error in sending request to %s: %s", url, err.Error()),
		})
	}
	defer resp.Body.Close()

	var rpcErr json2.Error
	err = json2.DecodeClientResponse(resp.Body, &rpcErr)
	jErr, ok := err.(*json2.Error)
	if ok {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: jErr.Error(),
		})
	}

	rerr := uc.repo.ActiveUser(id)
	if rerr != nil {
		return c.JSON(rerr.StatusCode(), response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		})
	}
	return c.JSON(http.StatusOK, response.Success{
		Code:    http.StatusOK,
		Message: "User Activated",
	})
}
func (uc *userController) DisableUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		}
		return c.JSON(rErr.Code, rErr)
	}

	user, rErr := uc.repo.FindUserByID(id)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	server, rErr := uc.repo.FindServerByID(user.ServerID)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
	params := &request.DisableUserRPC{
		Username: user.Username,
	}

	msg, err := json2.EncodeClientRequest("disableUser", params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Error in encoding request",
		})
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		})
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("Error in sending request to %s: %s", url, err.Error()),
		})
	}
	defer resp.Body.Close()

	var rpcErr json2.Error
	err = json2.DecodeClientResponse(resp.Body, &rpcErr)
	jErr, ok := err.(*json2.Error)
	if ok {
		return c.JSON(http.StatusInternalServerError, response.Error{
			Code:    http.StatusInternalServerError,
			Message: jErr.Error(),
		})
	}

	rerr := uc.repo.DisableUser(id)
	if rerr != nil {
		return c.JSON(rerr.StatusCode(), response.Error{
			Code:    rerr.StatusCode(),
			Message: rerr.Error(),
			Errors:  rerr.Errors(),
		})
	}
	return c.JSON(http.StatusOK, response.Success{
		Code:    http.StatusOK,
		Message: "User Disabled",
	})
}
