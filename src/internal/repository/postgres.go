package repository

import (
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/pkg/rest_err"
)

type Postgres interface {
	FindAdminByUsernameAndPAssword(username, password string) (*model.Admin, rest_err.RestErr)
	FindAllUsers() ([]model.User, rest_err.RestErr)
	FindUserByID(id int) (*model.User, rest_err.RestErr)
	CreateUser(user model.User) rest_err.RestErr
	UpdateUser(user model.User) rest_err.RestErr
	DeleteUser(id int) rest_err.RestErr
}
