package repository

import (
	"github.com/yazdsr/master-api/internal/entity/model"
	"github.com/yazdsr/master-api/internal/transport/http/request"
	"github.com/yazdsr/master-api/pkg/rest_err"
)

type Postgres interface {
	FindAdminByUsernameAndPAssword(username, password string) (*model.Admin, rest_err.RestErr)
	FindAllUsers() ([]model.User, rest_err.RestErr)
	FindUserByID(id int) (*model.User, rest_err.RestErr)
	CreateUser(user request.CreateUser) rest_err.RestErr
	UpdateUser(id int, user request.UpdateUser) rest_err.RestErr
	DeleteUser(id int) rest_err.RestErr
	FindServerByID(id int) (*model.Server, rest_err.RestErr)
	FindUserByUsernameAndServerID(username string, serverID int) (*model.User, rest_err.RestErr)
}
