package department

import (
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewDeportamentHandler(router *http.ServeMux, deps DepartmenthandlerDeps) {
	handler := &Departmenthandler{
		Config: deps.Config,
	}
	router.HandleFunc("/department/getDepartment", handler.getDepartment())

}

func (handler *Departmenthandler) getDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := req.HandleBody[DepartmentRequest](&w, r)
		if err != nil {
			return
		}

		data1 := DepartmentResponse{
			Id:   token.CreateId(),
			Name: token.CreateId(),
		}
		data2 := DepartmentResponse{
			Id:   token.CreateId(),
			Name: token.CreateId(),
		}
		data3 := DepartmentResponse{
			Id:   token.CreateId(),
			Name: token.CreateId(),
		}
		data4 := DepartmentResponse{
			Id:   token.CreateId(),
			Name: token.CreateId(),
		}
		var data []DepartmentResponse
		data = append(data, data1)
		data = append(data, data2)
		data = append(data, data3)
		data = append(data, data4)
		res.Json(w, data, 200)
	}
}
