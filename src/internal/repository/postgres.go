package repository

import (
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/pkg/rest_err"
)

type Postgres interface {
	FindAdminByUsernameAndPAssword(username, password string) (*model.Admin, rest_err.RestErr)
}
