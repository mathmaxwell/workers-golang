package department

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/pkg/db"
)

type Department struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type DeportamentRepository struct {
	DataBase *db.Db
}
type DepartmenthandlerDeps struct {
	*configs.Config
	DeportamentRepository *DeportamentRepository
}
type Departmenthandler struct {
	*configs.Config
	DeportamentRepository *DeportamentRepository
}
type DepartmentRequest struct {
	Name string `json:"name"`
}
type getDepartmentRequest struct {
	Token string `json:"token"`
}
