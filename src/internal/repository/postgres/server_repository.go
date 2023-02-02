package postgres

import (
	"context"
	"net/http"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/pkg/rest_err"
)

func (psql *postgres) FindServerByID(id int) (*model.Server, rest_err.RestErr) {
	server := new(model.Server)
	query := `SELECT * FROM servers WHERE id = $1`
	err := pgxscan.Get(context.Background(), psql.db, server, query, id)
	if err != nil {
		if err.Error() == "scanning one: no rows in result set" {
			return nil, rest_err.NewRestErr(http.StatusNotFound, "server not found", []string{})
		}
		return nil, rest_err.NewRestErr(http.StatusInternalServerError, err.Error(), []string{})
	}
	return server, nil
}
