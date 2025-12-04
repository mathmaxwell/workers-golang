package department

import "demo/purpleSchool/configs"

type DepartmenthandlerDeps struct {
	*configs.Config
}
type Departmenthandler struct {
	*configs.Config
}
type DepartmentRequest struct {
	Token string `json:"token" validate:"required"`
}
type DepartmentResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
