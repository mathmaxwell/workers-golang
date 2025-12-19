package employees

import (
	workschedule "demo/purpleSchool/internal/workSchedule"
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/fields"
	"demo/purpleSchool/pkg/files"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"gorm.io/gorm"
)

func NewEmployeeRepository(dataBase *db.Db) *EmployeeRepository {
	return &EmployeeRepository{
		DataBase: dataBase,
	}
}

func NewEmployeeHandler(router *http.ServeMux, deps EmployeeshandlerDeps) *EmployeesHandler {
	handler := &EmployeesHandler{
		Config:             deps.Config,
		EmployeeRepository: *deps.EmployeeRepository,
		AuthHandler:        deps.AuthHandler,
	}
	router.HandleFunc("/employees/createEmployees", handler.createEmployees())
	router.HandleFunc("/employees/updateEmployees", handler.updateEmployees())
	router.HandleFunc("/employees/getEmployees", handler.getEmployees())
	router.HandleFunc("/employees/deleteEmployee", handler.deleteEmployee())
	router.HandleFunc("/employees/getLateEmployeesById", handler.getLateEmployeesById())
	router.HandleFunc("/employees/getEmployeesById", handler.getEmployeesById())
	router.HandleFunc("/employees/getLateEmployees", handler.getLateEmployees())
	router.HandleFunc("/employees/getEmployeesCount", handler.getEmployeesCount())
	router.HandleFunc("/employees/getStatusById", handler.getStatusById())
	router.HandleFunc("/employees/createStatus", handler.createStatus())
	router.HandleFunc("/employees/getEmployeesByStatus", handler.getEmployeesByStatus())
	return handler
}
func (handler *EmployeesHandler) GetEmpById(id string) (Employee, error) {
	var employee Employee
	err := handler.EmployeeRepository.DataBase.Where("id = ?", id).First(&employee).Error
	if err != nil {
		return employee, errors.New("user is not found")
	}

	return employee, nil
}

func (handler *EmployeesHandler) createEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userToken := r.FormValue("token")
		user, err := handler.AuthHandler.GetUserByToken(userToken)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		var Accepted string = "false"
		if user.UserRole == 1 {
			Accepted = "true"
		}
		var token string = token.CreateId()
		if user.UserRole != 1 {
			token = userToken
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
		photoPath, err := files.SaveFile(file, header)
		if err != nil {
			res.Json(w, "failed to save image", http.StatusInternalServerError)
			return
		}

		newEmployee := Employee{
			Id:                         token,
			Gender:                     r.FormValue("gender"),
			Full_name:                  r.FormValue("full_name"),
			PINFL:                      r.FormValue("PINFL"),
			Phone_number:               r.FormValue("phone_number"),
			Passport_series_and_number: r.FormValue("passport_series_and_number"),
			Department:                 r.FormValue("department"),
			Position:                   r.FormValue("position"),
			Date_of_birth:              r.FormValue("date_of_birth"),
			Birth_month:                r.FormValue("birth_month"),
			Year_of_birth:              r.FormValue("year_of_birth"),
			Place_of_birth:             r.FormValue("place_of_birth"),
			Nationality:                r.FormValue("nationality"),
			Email:                      r.FormValue("Email"),
			Image:                      photoPath,
			Accepted:                   Accepted,
		}
		if err := fields.ValidateFields(newEmployee, Employee{}); err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		if err := handler.EmployeeRepository.DataBase.Create(&newEmployee).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, newEmployee, 200)
	}
}
func (handler *EmployeesHandler) getEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByIdRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		var employee Employee
		err = handler.EmployeeRepository.DataBase.Where("id = ?", body.Id).First(&employee).Error
		if err != nil {
			res.Json(w, "employee is not found", 400)
			return
		}
		res.Json(w, employee, 200)
	}
}
func (handler *EmployeesHandler) getEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		if body.Page <= 0 {
			body.Page = 1
		}
		if body.Count <= 0 {
			body.Count = 10
		}
		if body.Count > 100 {
			body.Count = 100
		}
		offset := (body.Page - 1) * body.Count
		query := handler.EmployeeRepository.DataBase.WithContext(r.Context())
		if body.Gender != "" {
			query = query.Where("gender ILIKE ?", "%"+body.Gender+"%")
		}
		if body.Passport_series_and_number != "" {
			query = query.Where("passport_series_and_number ILIKE ?", "%"+body.Passport_series_and_number+"%")
		}
		if body.PINFL != "" {
			query = query.Where("PINFL = ?", body.PINFL)
		}
		if body.Full_name != "" {
			query = query.Where("full_name ILIKE ?", "%"+body.Full_name+"%")
		}
		if body.Department != "" {
			query = query.Where("department ILIKE ?", "%"+body.Department+"%")
		}
		if body.Position != "" {
			query = query.Where("position ILIKE ?", "%"+body.Position+"%")
		}
		if body.Date_of_birth != "" {
			query = query.Where("date_of_birth = ?", body.Date_of_birth)
		}
		if body.Birth_month != "" {
			query = query.Where("birth_month = ?", body.Birth_month)
		}
		if body.Year_of_birth != "" {
			query = query.Where("year_of_birth = ?", body.Year_of_birth)
		}
		if body.Place_of_birth != "" {
			query = query.Where("place_of_birth ILIKE ?", "%"+body.Place_of_birth+"%")
		}
		if body.Nationality != "" {
			query = query.Where("nationality ILIKE ?", "%"+body.Nationality+"%")
		}
		allowedFields := map[string]string{
			"full_name":                  "full_name",
			"department":                 "department",
			"position":                   "position",
			"gender":                     "gender",
			"passport_series_and_number": "passport_series_and_number",
			"PINFL":                      "PINFL",
			"place_of_birth":             "place_of_birth",
			"nationality":                "nationality",
		}
		sortOrder := "ASC"
		if !body.SortAsc {
			sortOrder = "DESC"
		}
		if body.SortField == "date_of_birth" {
			orderClause := fmt.Sprintf(
				"MAKE_DATE(year_of_birth::int, birth_month::int, date_of_birth::int) %s",
				sortOrder,
			)
			query = query.Order(orderClause)
		} else {
			sortField := "full_name"
			if field, ok := allowedFields[body.SortField]; ok {
				sortField = field
			}
			query = query.Order(sortField + " " + sortOrder)
		}
		var employees []Employee
		err = query.
			Limit(body.Count).
			Offset(offset).
			Find(&employees).Error
		if err != nil {
			res.Json(w, "failed to get employees", 500)
			return
		}
		res.Json(w, employees, 200)
	}
}
func (handler *EmployeesHandler) deleteEmployee() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByIdRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		db := handler.EmployeeRepository.DataBase
		result := db.Delete(&Employee{}, "id = ?", body.Id)
		if result.Error != nil {
			res.Json(w, result.Error.Error(), 500)
			return
		}
		if result.RowsAffected == 0 {
			res.Json(w, "employee not found", 404)
			return
		}
		res.Json(w, "employee deleted", 200)
	}
}
func (handler *EmployeesHandler) getEmployeesCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[IEmployeesCountRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		targetDate := body.Day + body.Month*100 + body.Year*1000
		db := handler.EmployeeRepository.DataBase
		var onVacation, onSickLeave, onBusinessTrip, total int64
		db.Model(&EmployeeStatus{}).
			Where("status = ?", "on_vacation").
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Count(&onVacation)
		db.Model(&EmployeeStatus{}).
			Where("status = ?", "on_sick_leave").
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Count(&onSickLeave)
		db.Model(&EmployeeStatus{}).
			Where("status = ?", "on_a_business_trip").
			Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
			Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
			Count(&onBusinessTrip)
		var schedulesForDay []workschedule.ScheduleForDay
		db.Model(&workschedule.ScheduleForDay{}).Where("start_day + start_month*100 + start_year*1000 = ?", targetDate).Find(&schedulesForDay)
		absence := 0
		for _, value := range schedulesForDay {
			var count int64
			db.Model(&Record{}).
				Where("employee_id = ?", value.EmployeeId).
				Where("day::int + month::int*100 + year::int*1000 = ?", targetDate).
				Count(&count)
			if count == 0 {
				absence++
			}
		}
		db.Model(&Employee{}).Count(&total)
		data := IEmployeesCountResponse{
			On_vacation:        int(onVacation),
			On_sick_leave:      int(onSickLeave),
			On_a_business_trip: int(onBusinessTrip),
			Absence:            int(absence),
			Total_employees:    int(total),
		}
		res.Json(w, data, 200)
	}
}
func (handler *EmployeesHandler) updateEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			res.Json(w, "не удалось разобрать форму", http.StatusBadRequest)
			return
		}
		userToken := r.FormValue("token")
		if userToken == "" {
			res.Json(w, "токен обязателен", http.StatusUnauthorized)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(userToken)
		if err != nil {
			res.Json(w, "пользователь не найден или токен недействителен", http.StatusUnauthorized)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "доступ запрещён: требуется роль администратора", http.StatusForbidden)
			return
		}
		employeeId := r.FormValue("id")
		if employeeId == "" {
			res.Json(w, "ID сотрудника обязателен", http.StatusBadRequest)
			return
		}
		var employee Employee
		if err := handler.EmployeeRepository.DataBase.Where("id = ?", employeeId).First(&employee).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.Json(w, "сотрудник не найден", http.StatusNotFound)
			} else {
				res.Json(w, "ошибка базы данных", http.StatusInternalServerError)
			}
			return
		}
		photoPath := employee.Image
		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			defer file.Close()
			var saveErr error
			photoPath, saveErr = files.SaveFile(file, header)
			if saveErr != nil {
				res.Json(w, "не удалось сохранить изображение", http.StatusInternalServerError)
				return
			}
		}
		updatedEmployee := Employee{
			Gender:                     fields.GetOrDefault(r.FormValue("gender"), employee.Gender),
			Full_name:                  fields.GetOrDefault(r.FormValue("full_name"), employee.Full_name),
			PINFL:                      fields.GetOrDefault(r.FormValue("PINFL"), employee.PINFL),
			Phone_number:               fields.GetOrDefault(r.FormValue("phone_number"), employee.Phone_number),
			Passport_series_and_number: fields.GetOrDefault(r.FormValue("passport_series_and_number"), employee.Passport_series_and_number),
			Department:                 fields.GetOrDefault(r.FormValue("department"), employee.Department),
			Position:                   fields.GetOrDefault(r.FormValue("position"), employee.Position),
			Date_of_birth:              fields.GetOrDefault(r.FormValue("date_of_birth"), employee.Date_of_birth),
			Birth_month:                fields.GetOrDefault(r.FormValue("birth_month"), employee.Birth_month),
			Year_of_birth:              fields.GetOrDefault(r.FormValue("year_of_birth"), employee.Year_of_birth),
			Place_of_birth:             fields.GetOrDefault(r.FormValue("place_of_birth"), employee.Place_of_birth),
			Nationality:                fields.GetOrDefault(r.FormValue("nationality"), employee.Nationality),
			Email:                      fields.GetOrDefault(r.FormValue("Email"), employee.Email),
			Image:                      photoPath,
			Accepted:                   fields.GetOrDefault(r.FormValue("accepted"), employee.Accepted),
		}
		if err := handler.EmployeeRepository.DataBase.Model(&employee).Updates(updatedEmployee).Error; err != nil {
			res.Json(w, "не удалось обновить сотрудника", http.StatusInternalServerError)
			return
		}
		res.Json(w, map[string]interface{}{
			"message":  "сотрудник успешно обновлён",
			"employee": updatedEmployee,
		}, http.StatusOK)
	}
}
func (handler *EmployeesHandler) createStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createStatusRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		data := EmployeeStatus{
			Id:         token.CreateId(),
			EmployeeId: body.EmployeeId,
			Status:     body.Status,
			StartDay:   body.StartDay,
			StartMonth: body.StartMonth,
			StartYear:  body.StartYear,
			EndDay:     body.EndDay,
			EndMonth:   body.EndMonth,
			EndYear:    body.EndYear,
		}
		if err := handler.EmployeeRepository.DataBase.Create(&data).Error; err != nil {
			res.Json(w, "db error", 500)
			return
		}
		res.Json(w, body, 200)
	}
}
func (handler *EmployeesHandler) getStatusById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByIdRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 401)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}

		db := handler.EmployeeRepository.DataBase
		var employeeStatusHistory []EmployeeStatus
		db.Model(&EmployeeStatus{}).
			Where("employee_id = ?", body.Id).
			Find(&employeeStatusHistory)
		res.Json(w, employeeStatusHistory, 200)
	}
}
func (handler *EmployeesHandler) getEmployeesByStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetEmployeesByStatusRequest](&w, r)
		if err != nil {
			res.Json(w, err.Error(), 401)
			return
		}
		user, err := handler.AuthHandler.GetUserByToken(body.Token)
		if err != nil {
			res.Json(w, "user is not found", 401)
			return
		}
		if user.UserRole != 1 {
			res.Json(w, "you are not admin", 403)
			return
		}
		db := handler.EmployeeRepository.DataBase
		targetDate := body.Day + body.Month*100 + body.Year*1000
		var ids []string
		if body.Status == "absence" {
			var schedulesForDay []workschedule.ScheduleForDay
			db.Model(&workschedule.ScheduleForDay{}).Where("start_day + start_month*100 + start_year*1000 = ?", targetDate).Find(&schedulesForDay)
			for _, value := range schedulesForDay {
				var count int64
				db.Model(&Record{}).
					Where("employee_id = ?", value.EmployeeId).
					Where("day::int + month::int*100 + year::int*1000 = ?", targetDate).
					Count(&count)
				if count == 0 {
					ids = append(ids, value.EmployeeId)
				}
			}
		} else {
			db.Model(&EmployeeStatus{}).Select("EmployeeId").
				Where("status = ?", body.Status).
				Where("start_day + start_month*100 + start_year*1000 <= ?", targetDate).
				Where("end_day + end_month*100 + end_year*1000 >= ?", targetDate).
				Find(&ids)
		}
		res.Json(w, ids, 200)
	}
}
func (handler *EmployeesHandler) getLateEmployeesById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getLateEmployeesByIdRequest](&w, r)
		if err != nil {
			res.Json(w, err, 400)
			return
		}
		targetStart := body.Start_day + body.Start_month*100 + body.Start_year*10000
		targetEnd := body.End_day + body.End_month*100 + body.End_year*10000
		db := handler.EmployeeRepository.DataBase
		var schedules []workschedule.ScheduleForDay
		err = db.Model(&workschedule.ScheduleForDay{}).
			Where("employee_id = ? AND start_day + start_month*100 + start_year*10000 >= ? AND start_day + start_month*100 + start_year*10000 <= ?",
				body.EmployeeId, targetStart, targetEnd).
			Find(&schedules).Error
		if err != nil {
			res.Json(w, err.Error(), 500)
			return
		}
		result := make([]ITardinessHistory, 0, len(schedules))
		for _, schedule := range schedules {
			dayStart := schedule.StartDay + schedule.StartMonth*100 + schedule.StartYear*10000
			dayEnd := schedule.EndDay + schedule.EndMonth*100 + schedule.EndYear*10000
			var records []Record
			err = db.Model(&Record{}).
				Where("employee_id = ?", schedule.EmployeeId).
				Where("day::int + month::int*100 + year::int*10000 >= ? AND day::int + month::int*100 + year::int*10000 <= ?",
					dayStart, dayEnd).
				Find(&records).Error
			if err != nil {
				res.Json(w, err.Error(), 500)
				return
			}
			sort.Slice(records, func(i, j int) bool {
				dateI := toInt(records[i].Year)*10000 + toInt(records[i].Month)*100 + toInt(records[i].Day)
				dateJ := toInt(records[j].Year)*10000 + toInt(records[j].Month)*100 + toInt(records[j].Day)
				if dateI != dateJ {
					return dateI < dateJ
				}
				timeI := toInt(records[i].Hour)*3600 + toInt(records[i].Minute)*60 + toInt(records[i].Second)
				timeJ := toInt(records[j].Hour)*3600 + toInt(records[j].Minute)*60 + toInt(records[j].Second)
				return timeI < timeJ
			})
			history := ITardinessHistory{
				EmployeeId: body.EmployeeId,
				Year:       strconv.Itoa(schedule.StartYear),
				Month:      strconv.Itoa(schedule.StartMonth),
				Day:        strconv.Itoa(schedule.StartDay),
				WorkSchedule: ScheduleForDay{
					Id:         schedule.Id,
					EmployeeId: schedule.EmployeeId,
					StartHour:  schedule.StartHour,
					StartDay:   schedule.StartDay,
					StartMonth: schedule.StartMonth,
					StartYear:  schedule.StartYear,
					EndHour:    schedule.EndHour,
					EndDay:     schedule.EndDay,
					EndMonth:   schedule.EndMonth,
					EndYear:    schedule.EndYear,
				},
			}
			if len(records) == 0 {
				history.EntryHour = 99
				history.EntryMinute = 99
				history.EntryDay = schedule.StartDay
				history.EntryMonth = schedule.StartMonth
				history.EntryYear = schedule.StartYear

				history.ExitHour = 99
				history.ExitMinute = 99
				history.ExitDay = schedule.EndDay
				history.ExitMonth = schedule.EndMonth
				history.ExitYear = schedule.EndYear
			} else {
				first := records[0]
				history.EntryHour = toInt(first.Hour)
				history.EntryMinute = toInt(first.Minute)
				history.EntryDay = toInt(first.Day)
				history.EntryMonth = toInt(first.Month)
				history.EntryYear = toInt(first.Year)
				last := records[len(records)-1]
				history.ExitHour = toInt(last.Hour)
				history.ExitMinute = toInt(last.Minute)
				history.ExitDay = toInt(last.Day)
				history.ExitMonth = toInt(last.Month)
				history.ExitYear = toInt(last.Year)
			}
			result = append(result, history)
		}
		res.Json(w, result, 200)
	}
}

func (handler *EmployeesHandler) getLateEmployees() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[GetLateEmployeesRequest](&w, r)
		if body.Token == "" {
			return
		} //проверка токена
		if err != nil {
			return
		}
		//ищет за тех дней кто опоздал и сколько опоздал, график работы
		var data []workschedule.ITardinessHistory
		res.Json(w, data, 200)
	}
}
