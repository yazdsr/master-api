package postgres

import (
	"context"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/internal/transport/http/request"
	"github.com/yazdsr/master-api/pkg/hash"
	"github.com/yazdsr/master-api/pkg/rest_err"
)

func (psql *postgres) FindAllUsers() ([]model.User, rest_err.RestErr) {
	var users []model.User = []model.User{}
	query := `SELECT * FROM users ORDER BY id ASC`
	err := pgxscan.Select(context.Background(), psql.db, &users, query)
	if err != nil {
		return []model.User{}, rest_err.NewRestErr(http.StatusInternalServerError, err.Error(), []string{})
	}
	return users, nil
}

func (psql *postgres) FindUserByID(id int) (*model.User, rest_err.RestErr) {
	user := new(model.User)
	query := `SELECT * FROM users WHERE id = $1`
	err := pgxscan.Get(context.Background(), psql.db, user, query, id)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return nil, rest_err.NewRestErr(http.StatusNotFound, "user not found", []string{})
		}
		return nil, rest_err.NewRestErr(http.StatusInternalServerError, err.Error(), []string{})
	}
	return user, nil
}

func (psql *postgres) CreateUser(usr request.CreateUser) rest_err.RestErr {
	user := new(model.User)
	password := hash.GenerateSha256(usr.Password)
	query := `INSERT INTO users (username, password, full_name, server_id, valid_until) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := psql.db.QueryRow(context.Background(), query, usr.Username, password, usr.FullName, usr.ServerID, usr.ValidUntil).Scan(&user.ID)
	if err != nil {
		return rest_err.NewRestErr(http.StatusInternalServerError, "error while adding user", []string{err.Error()})
	}
	if user.ID == 0 {
		return rest_err.NewRestErr(http.StatusInternalServerError, "error while adding user", []string{})
	}
	return nil
}

func (psql *postgres) UpdateUser(id int, user request.UpdateUser) rest_err.RestErr {
	var query string
	password := hash.GenerateSha256(user.Password)
	var err error
	if user.Password == "" {
		query = `UPDATE users SET full_name = $1, valid_until = $2 WHERE id = $3`
		_, err = psql.db.Exec(context.Background(), query, user.FullName, user.ValidUntil, id)
	} else {
		query = `UPDATE users SET password = $1, full_name = $2, valid_until = $3 WHERE id = $4`
		_, err = psql.db.Exec(context.Background(), query, password, user.FullName, user.ValidUntil, id)
	}
	if err != nil {
		return rest_err.NewRestErr(http.StatusInternalServerError, "error while updating user", []string{err.Error()})
	}
	return nil
}

func (psql *postgres) ActiveUser(id int) rest_err.RestErr {
	var query string
	var err error

	query = `UPDATE users SET active = $1 WHERE id = $2`
	_, err = psql.db.Exec(context.Background(), query, true, id)

	if err != nil {
		return rest_err.NewRestErr(http.StatusInternalServerError, "error while activating user", []string{err.Error()})
	}
	return nil
}

func (psql *postgres) DisableUser(id int) rest_err.RestErr {
	var query string
	var err error

	query = `UPDATE users SET active = $1 WHERE id = $2`
	_, err = psql.db.Exec(context.Background(), query, false, id)

	if err != nil {
		return rest_err.NewRestErr(http.StatusInternalServerError, "error while activating user", []string{err.Error()})
	}
	return nil
}

func (psql *postgres) DeleteUser(id int) rest_err.RestErr {
	query := `DELETE FROM users WHERE id = $1`
	_, err := psql.db.Exec(context.Background(), query, id)
	if err != nil {
		return rest_err.NewRestErr(http.StatusInternalServerError, "error while deleting user", []string{err.Error()})
	}
	return nil
}

func (psql *postgres) FindUserByUsernameAndServerID(username string, serverID int) (*model.User, rest_err.RestErr) {
	user := new(model.User)
	query := `SELECT * FROM users WHERE username = $1 AND server_id = $2`
	err := pgxscan.Get(context.Background(), psql.db, user, query, username, serverID)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return nil, rest_err.NewRestErr(http.StatusNotFound, "user not found", []string{})
		}
		return nil, rest_err.NewRestErr(http.StatusInternalServerError, err.Error(), []string{})
	}
	return user, nil
}
