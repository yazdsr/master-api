package postgres

import (
	"context"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/pkg/hash"
	"github.com/yazdsr/master-api/pkg/rest_err"
)

func (psql *postgres) FindAdminByUsernameAndPAssword(username, password string) (*model.Admin, rest_err.RestErr) {
	admin := new(model.Admin)
	password = hash.GenerateSha256(password)
	query := `SELECT id, username, password, created_at, updated_at FROM admins WHERE username = $1 AND password = $2`
	if err := pgxscan.Get(context.Background(), psql.db, admin, query, username, password); err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return nil, rest_err.NewRestErr(http.StatusUnauthorized, "invalid username or password", []string{})
		}
		return nil, rest_err.NewRestErr(http.StatusInternalServerError, err.Error(), []string{})
	}
	return admin, nil
}
