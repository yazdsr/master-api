package echo

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/rpc/v2/json2"
	"github.com/labstack/echo/v4"
	"github.com/yazdsr/master-api/internal/pkg/logger"
	"github.com/yazdsr/master-api/internal/repository"
	"github.com/yazdsr/master-api/internal/transport/http/response"
)

type serverController struct {
	logger logger.Logger
	repo   repository.Postgres
}

func (sc *serverController) HeartBeat(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rErr := response.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID",
		}
		return c.JSON(rErr.Code, rErr)
	}

	server, rErr := sc.repo.FindServerByID(id)
	if rErr != nil {
		return c.JSON(rErr.StatusCode(), response.Error{
			Code:    rErr.StatusCode(),
			Message: rErr.Error(),
		})
	}

	errChan := make(chan error)

	ctxTimeout, cancel := context.WithTimeout(c.Request().Context(), time.Second*2)
	defer cancel()

	go func(errChan chan error) {
		url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
		msg, err := json2.EncodeClientRequest("heartbeat", struct{}{})
		if err != nil {
			errChan <- c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: "Error in encoding request",
			})
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
		if err != nil {
			errChan <- c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			})
		}

		req.Header.Set("Content-Type", "application/json")
		client := new(http.Client)
		resp, err := client.Do(req)
		if err != nil {
			errChan <- c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("Error in sending request to %s: %s", url, err.Error()),
			})
		}
		defer resp.Body.Close()

		var result string
		err = json2.DecodeClientResponse(resp.Body, &result)
		if err != nil {
			errChan <- c.JSON(http.StatusInternalServerError, response.Error{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
		}
		errChan <- c.JSON(http.StatusOK, response.Success{
			Code:    http.StatusOK,
			Message: result,
		})
	}(errChan)

	select {
	case <- ctxTimeout.Done():
		return c.JSON(http.StatusRequestTimeout, response.Error{
			Code: http.StatusRequestTimeout,
			Message: "Server Isn't Alive!",
		})
	case err = <- errChan:
		return err
	}
}
