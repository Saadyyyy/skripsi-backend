package route

import (
	categoryhandler "bank_soal/api/category/category_handler"
	categoryrepository "bank_soal/api/category/category_repository"
	categoryservice "bank_soal/api/category/category_service"
	"bank_soal/api/soal/soal_handler"
	"bank_soal/api/soal/soal_repository"
	"bank_soal/api/soal/soal_service"
	handlerUser "bank_soal/api/user/user_handler"
	repositoryUser "bank_soal/api/user/user_repository"
	serviceUser "bank_soal/api/user/user_service"
	"bank_soal/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Register(db *sqlx.DB, echo *echo.Echo) {

	//user
	repositoryUser := repositoryUser.NewUserRepository(db)
	serviceUser := serviceUser.NewUserService(repositoryUser, db)
	handlerUser := handlerUser.NewUserHandler(serviceUser)
	//route
	user := echo.Group("/user")
	user.POST("/register", handlerUser.CreateUser, middleware.JWTMiddleware(), middleware.AdminMiddleware).Name = "CreateUser"
	user.GET("/", handlerUser.GetAllUser, middleware.JWTMiddleware(), middleware.AdminMiddleware).Name = "GetAllUser"
	user.POST("/login", handlerUser.LoginUser).Name = "LoginUser"
	user.POST("/update", handlerUser.UpdateUser, middleware.JWTMiddleware(), middleware.AdminMiddleware).Name = "UpdateUser"

	//soal
	repositorySoal := soal_repository.NewSoalRepository(db)
	serviceSoal := soal_service.NewSoalService(repositorySoal, db)
	handlersoal := soal_handler.NewSoalHandler(serviceSoal)
	//route
	soal := echo.Group("/soal", middleware.JWTMiddleware())
	soal.POST("/create", handlersoal.CreateSoal, middleware.AdminMiddleware).Name = "CreateSoal"
	soal.GET("/", handlersoal.GetSoal).Name = "GetSoal"
	soal.POST("/update", handlersoal.UpdateSoal).Name = "UpdateSoal"
	soal.POST("/delete", handlersoal.DeletedSoal).Name = "DeletedSoal"
	soal.GET("/detail", handlersoal.GetSoalById).Name = "GetSoalById"

	//category
	repositoryCategory := categoryrepository.NewCategoryRepository(db)
	serviceCategory := categoryservice.NewCategoryService(repositoryCategory)
	handlerCategory := categoryhandler.NewCategoryHandler(serviceCategory)
	//route
	category := echo.Group("/category", middleware.JWTMiddleware())
	category.POST("/create", handlerCategory.CreateSoal).Name = "CreateSoal"
	category.GET("/detail", handlerCategory.GetCategoryByID).Name = "GetCategoryByID"
	category.GET("/", handlerCategory.GetListCategory).Name = "GetListCategory"
}
