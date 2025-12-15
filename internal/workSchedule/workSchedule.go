package workschedule

import (
	"demo/purpleSchool/pkg/db"
	"demo/purpleSchool/pkg/req"
	"demo/purpleSchool/pkg/res"
	"demo/purpleSchool/pkg/token"
	"net/http"
)

func NewScheduleRepository(dataBase *db.Db) *ScheduleRepository {
	return &ScheduleRepository{
		DataBase: dataBase,
	}
}

func NewWorkScheduleHandler(router *http.ServeMux, deps WorkScheduleDeps) {
	handler := &WorkScheduleHandler{
		Config:             deps.Config,
		ScheduleRepository: deps.ScheduleRepository,
		AuthHandler:        deps.AuthHandler,
	}
	router.HandleFunc("/workSchedule/getEmployeeWorkSchedule", handler.getEmployeeWorkSchedule())
	router.HandleFunc("/workSchedule/createWorkSchedule", handler.createWorkSchedule())
	router.HandleFunc("/workSchedule/updateWorkSchedule", handler.updateWorkSchedule())
}
func (handler *WorkScheduleHandler) createWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[createWorkScheduleRequest](&w, r)
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
		newSchedule := ScheduleForDay{
			Id:         token.CreateId(),
			EmployeeId: body.EmployeeId,
			StartHour:  body.StartHour,
			StartDay:   body.StartDay,
			StartMonth: body.StartMonth,
			StartYear:  body.StartYear,
			EndHour:    body.EndHour,
			EndDay:     body.EndDay,
			EndMonth:   body.EndMonth,
			EndYear:    body.EndYear,
		}
		err = handler.ScheduleRepository.DataBase.Create(&newSchedule).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}
		res.Json(w, newSchedule, 200)
	}
}

func (handler *WorkScheduleHandler) getEmployeeWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[getWorkScheduleRequest](&w, r)
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
		schedule := []ScheduleForDay{}
		startDate := body.StartYearSchedule*10000 +
			body.StartMonthSchedule*100 +
			body.StartDaySchedule
		endDate := body.EndYearSchedule*10000 +
			body.EndMonthSchedule*100 +
			body.EndDaySchedule
		db := handler.ScheduleRepository.DataBase
		db.Model(&ScheduleForDay{}).
			Where(
				"(start_year*10000 + start_month*100 + start_day) <= ? AND "+
					"(end_year*10000 + end_month*100 + end_day) >= ?",
				endDate,
				startDate,
			).Where("employee_id = ?", body.EmployeeId).
			Find(&schedule)
		response := getWorkScheduleResponse{
			StartDaySchedule:   body.StartDaySchedule,
			StartMonthSchedule: body.StartMonthSchedule,
			StartYearSchedule:  body.StartYearSchedule,
			EndDaySchedule:     body.EndDaySchedule,
			EndMonthSchedule:   body.EndMonthSchedule,
			EndYearSchedule:    body.EndYearSchedule,
			WorkSchedule:       schedule,
		}
		res.Json(w, response, 200)
	}
}

func (handler *WorkScheduleHandler) updateWorkSchedule() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[updateWorkScheduleRequest](&w, r)
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
		db := handler.ScheduleRepository.DataBase
		var schedule ScheduleForDay
		err = db.Where("id = ?", body.Id).First(&schedule).Error
		if err != nil {
			res.Json(w, err.Error(), 400)
		}
		db.Model(&schedule).Updates(body.ScheduleForDay)
		res.Json(w, body, 200)
	}
}
