package records

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/fields"
	"demo/purpleSchool/pkg/files"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
	"strconv"
	"strings"
)

func NewRecordRepository(dataBase *db.Db) *RecordRepository {
	return &RecordRepository{
		DataBase: dataBase,
	}
}
func NewRecordHandler(router *http.ServeMux, deps RecordhandlerDeps) *RecordHandler {
	handler := &RecordHandler{
		Config:             deps.Config,
		RecordRepository:   *deps.RecordRepository,
		EmployeeRepository: deps.EmployeeRepository,
	}
	router.HandleFunc("/record/create", handler.create())
	return handler
}
func (handler *RecordHandler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := handler.EmployeeRepository.GetEmpById(r.FormValue("employeeId"))
		if err != nil {
			res.Json(w, "user is not found for record", http.StatusBadRequest)
			return
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			res.Json(w, "failed to parse form", http.StatusBadRequest)
			return
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			res.Json(w, "image is not found", 400)
			return
		}
		photoPath, err := files.SaveRecord(file, header)
		if err != nil {
			res.Json(w, "failed to save image", http.StatusInternalServerError)
			return
		}
		year, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("year")))
		month, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("month")))
		day, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("day")))
		hour, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("hour")))
		minute, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("minute")))
		second, _ := strconv.Atoi(strings.TrimSpace(r.FormValue("second")))
		newRecord := Record{
			Id:          token.CreateId(),
			EmployeeId:  r.FormValue("employeeId"),
			Image:       photoPath,
			Year:        year,
			Month:       month,
			Day:         day,
			Hour:        hour,
			Minute:      minute,
			Second:      second,
			Description: r.FormValue("description"),
		}
		if err := fields.ValidateFields(newRecord, Record{}); err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		if err := handler.RecordRepository.DataBase.Create(&newRecord).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, newRecord, 200)
	}
}
