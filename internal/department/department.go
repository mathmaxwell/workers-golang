package department

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewDeportamentRepository(dataBase *db.Db) *DeportamentRepository {
	return &DeportamentRepository{
		DataBase: dataBase,
	}
}
func NewDeportamentHandler(router *http.ServeMux, deps DepartmenthandlerDeps) {
	handler := &Departmenthandler{
		Config:                deps.Config,
		DeportamentRepository: deps.DeportamentRepository,
	}
	router.HandleFunc("/department/createDepartment", handler.createDepartment())
	router.HandleFunc("/department/getDepartment", handler.getDepartment())
}
func (handler *Departmenthandler) createDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[DepartmentRequest](&w, r)
		if err != nil {
			return
		}
		newDeportament := Department{
			Name: body.Name,
			Id:   token.CreateId(),
		}
		handler.DeportamentRepository.DataBase.Create(&newDeportament)
		res.Json(w, newDeportament, 200)
	}
}

func (handler *Departmenthandler) getDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := req.HandleBody[getDepartmentRequest](&w, r)
		if err != nil {
			return
		}
		var data []Department
		err = handler.DeportamentRepository.DataBase.Find(&data).Error
		if err != nil {
			res.Json(w, err, 200)
			return
		}
		res.Json(w, data, 200)
	}
}
