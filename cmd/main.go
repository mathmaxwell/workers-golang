package main

import (
	"demo/purpleSchool/configs"
	"demo/purpleSchool/internal/auth"
	"demo/purpleSchool/internal/department"
	"demo/purpleSchool/internal/employees"
	"demo/purpleSchool/internal/messages"
	"demo/purpleSchool/internal/workSchedule"
	"demo/purpleSchool/pkg/cors"
	"demo/purpleSchool/pkg/db"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDB(conf)
	//repositories
	router := http.NewServeMux()
	mux := cors.Cors(router)

	//register
	authRepo := auth.NewUserRepository(database)
	database.AutoMigrate(&auth.User{})
	authHandler := auth.NewAuthHandler(router, auth.AuthhandlerDeps{
		Config:         conf,
		AuthRepository: authRepo,
	})

	//employees
	employeeRepo := employees.NewEmployeeRepository(database)
	database.AutoMigrate(&employees.IEmployeesResponse{})
	employees.NewEmployeeHandler(router, employees.EmployeeshandlerDeps{
		Config:             conf,
		EmployeeRepository: employeeRepo,
		AuthHandler:        authHandler,
	})

	//work schedule
	workschedule.NewWorkScheduleHandler(router, workschedule.WorkScheduleDeps{
		Config: conf,
	})
	//message
	messages.NewMessagesHandler(router, messages.MessagehandlerDeps{
		Config: conf,
	})
	//deportament
	department.NewDeportamentHandler(router, department.DepartmenthandlerDeps{
		Config: conf,
	})
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
