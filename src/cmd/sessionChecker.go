package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
	"github.com/yazdsr/master-api/internal/config"
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/internal/transport/http/request"
)

var sessionCheckerCMD = &cli.Command{
	Name:    "session_checker",
	Aliases: []string{"c"},
	Usage:   "check user session and disable expired sessions",
	Action:  sessionChecker,
}

func sessionChecker(c *cli.Context) error {
	cfg := new(config.Config)
	config.ReadEnv(cfg)
	dbCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db, err := pgxpool.Connect(dbCtx, url(cfg.Psql))
	if err != nil {
		return err
	}

	defer db.Close()

	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			fmt.Println("Fetching users...")
			var users []model.User = []model.User{}
			query := `SELECT * FROM users ORDER BY id ASC`
			err := pgxscan.Select(context.Background(), db, &users, query)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			for _, user := range users {
				if user.ValidUntil.Unix() < time.Now().Unix() || time.Now().Unix() < user.StartDate.Unix() {
					server := new(model.Server)
					query := `SELECT * FROM servers WHERE id = $1`
					err := pgxscan.Get(context.Background(), db, server, query, user.ServerID)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
					params := &request.DisableUserRPC{
						Username: user.Username,
					}

					msg, err := json2.EncodeClientRequest("disableUser", params)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					req.Header.Set("Content-Type", "application/json")
					client := new(http.Client)
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					var rpcErr json2.Error
					err = json2.DecodeClientResponse(resp.Body, &rpcErr)
					jErr, ok := err.(*json2.Error)
					if ok {
						fmt.Println(jErr.Error())
						continue
					}
					resp.Body.Close()

					query = `UPDATE users SET active = $1 WHERE id = $2`
					_, err = db.Exec(context.Background(), query, false, user.ID)

					if err != nil {
						fmt.Println(jErr.Error())
						continue
					}
				} else {
					server := new(model.Server)
					query := `SELECT * FROM servers WHERE id = $1`
					err := pgxscan.Get(context.Background(), db, server, query, user.ServerID)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					url := fmt.Sprintf("http://%s:%d", server.Ip, server.Port)
					params := &request.DisableUserRPC{
						Username: user.Username,
					}

					msg, err := json2.EncodeClientRequest("activeUser", params)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg))
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					req.Header.Set("Content-Type", "application/json")
					client := new(http.Client)
					resp, err := client.Do(req)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}

					var rpcErr json2.Error
					err = json2.DecodeClientResponse(resp.Body, &rpcErr)
					jErr, ok := err.(*json2.Error)
					if ok {
						fmt.Println(jErr.Error())
						continue
					}
					resp.Body.Close()

					query = `UPDATE users SET active = $1 WHERE id = $2`
					_, err = db.Exec(context.Background(), query, true, user.ID)

					if err != nil {
						fmt.Println(jErr.Error())
						continue
					}
				}
			}
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	fmt.Println("\nReceived an interrupt, closing connections...")

	return nil
}

func url(cfg config.Psql) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
}
