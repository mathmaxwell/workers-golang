package records

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/employees"
	"demo/purpleSchool/pkg/db"
)

type RecordRepository struct {
	DataBase *db.Db
}
type RecordhandlerDeps struct {
	*configs.Config
	RecordRepository   *RecordRepository
	EmployeeRepository *employees.EmployeesHandler
}
type RecordHandler struct {
	*configs.Config
	RecordRepository   RecordRepository
	EmployeeRepository *employees.EmployeesHandler
}
type Record struct {
	Id          string `json:"id"`
	EmployeeId  string `json:"employeeId"`
	Image       string `json:"image"`
	Year        int    `json:"year"`
	Month       int    `json:"month"`
	Day         int    `json:"day"`
	Hour        int    `json:"hour"`
	Minute      int    `json:"minute"`
	Second      int    `json:"second"`
	Description string `json:"description"`
}
