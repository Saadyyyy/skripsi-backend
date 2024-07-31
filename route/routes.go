package route

import (
	"bank_soal/api/soal/soal_handler"
	"bank_soal/api/soal/soal_repository"
	"bank_soal/api/soal/soal_service"
	handlerUser "bank_soal/api/user/user_handler"
	repositoryUser "bank_soal/api/user/user_repository"
	serviceUser "bank_soal/api/user/user_service"
	"bank_soal/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func Register(db *sqlx.DB, echo *echo.Echo) {

	//user
	repositoryUser := repositoryUser.NewUserRepository(db)
	serviceUser := serviceUser.NewUserService(repositoryUser, db)
	handlerUser := handlerUser.NewUserHandler(serviceUser)

	user := echo.Group("/user")
	user.GET("/", handlerUser.GetAllUser, middleware.JWTMiddleware()).Name = "GetAllUser"
	user.POST("/register", handlerUser.CreateUser).Name = "CreateUser"
	user.POST("/login", handlerUser.LoginUser).Name = "LoginUser"
	user.POST("/update", handlerUser.UpdateUser, middleware.JWTMiddleware()).Name = "UpdateUser"

	//soal
	repositorySoal := soal_repository.NewSoalRepository(db)
	serviceSoal := soal_service.NewSoalService(repositorySoal, db)
	handlersoal := soal_handler.NewSoalHandler(serviceSoal)

	soal := echo.Group("/soal", middleware.JWTMiddleware())
	soal.POST("/create", handlersoal.CreateSoal).Name = "CreateSoal"
	soal.GET("/", handlersoal.GetSoal).Name = "GetSoal"
	soal.POST("/update", handlersoal.UpdateSoal).Name = "UpdateSoal"
	soal.POST("/delete", handlersoal.DeletedSoal).Name = "DeletedSoal"
	soal.GET("/detail", handlersoal.GetSoalById).Name = "GetSoalById"
}
