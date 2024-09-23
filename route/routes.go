package route

import (
	categoryhandler "bank_soal/api/category/category_handler"
	categoryrepository "bank_soal/api/category/category_repository"
	categoryservice "bank_soal/api/category/category_service"
	rangkinghandler "bank_soal/api/rangking/rangking_handler"
	rangkingrepository "bank_soal/api/rangking/rangking_repository"
	rangkingservice "bank_soal/api/rangking/rangking_service"
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
	user.POST("/register", handlerUser.CreateUser).Name = "CreateUser"
	user.GET("/", handlerUser.GetAllUser).Name = "GetAllUser"
	// user.GET("/", handlerUser.GetAllUser, middleware.JWTMiddleware(), middleware.AdminMiddleware).Name = "GetAllUser"
	user.POST("/login", handlerUser.LoginUser).Name = "LoginUser"
	user.POST("/update", handlerUser.UpdateUser, middleware.JWTMiddleware(), middleware.AdminMiddleware).Name = "UpdateUser"
	user.GET("/detail", handlerUser.GetUserByID, middleware.JWTMiddleware()).Name = "GetUserByID"
	user.POST("/role", handlerUser.UpdateUserRoleByID, middleware.JWTMiddleware()).Name = "UpdateUserRoleByID"

	//soal
	repositorySoal := soal_repository.NewSoalRepository(db)
	serviceSoal := soal_service.NewSoalService(repositorySoal, db)
	handlersoal := soal_handler.NewSoalHandler(serviceSoal)
	//route
	soal := echo.Group("/soal")
	soal.POST("/create", handlersoal.CreateSoal, middleware.JWTMiddleware()).Name = "CreateSoal"
	soal.GET("/", handlersoal.GetSoal, middleware.JWTMiddleware()).Name = "GetSoal"
	soal.POST("/update", handlersoal.UpdateSoal, middleware.JWTMiddleware()).Name = "UpdateSoal"
	soal.POST("/delete", handlersoal.DeletedSoal, middleware.JWTMiddleware()).Name = "DeletedSoal"
	soal.GET("/detail", handlersoal.GetSoalById, middleware.JWTMiddleware()).Name = "GetSoalById"

	//category
	repositoryCategory := categoryrepository.NewCategoryRepository(db)
	serviceCategory := categoryservice.NewCategoryService(repositoryCategory)
	handlerCategory := categoryhandler.NewCategoryHandler(serviceCategory)
	//route
	category := echo.Group("/category")
	category.POST("/create", handlerCategory.CreateSoal, middleware.JWTMiddleware(), middleware.AdminMiddleware).Name = "CreateSoal"
	category.GET("/detail", handlerCategory.GetCategoryByID).Name = "GetCategoryByID"
	category.GET("/", handlerCategory.GetListCategory).Name = "GetListCategory"
	category.POST("/update", handlerCategory.UpdatedCategory).Name = "UpdatedCategory"
	category.POST("/delete", handlerCategory.DeletedCategory).Name = "DeletedSoal"

	// rangkings
	repoRank := rangkingrepository.NewRangkingRepository(db)
	serviceRank := rangkingservice.NewRangkingService(repoRank, repositorySoal, repositoryCategory, repositoryUser)
	handlerRank := rangkinghandler.NewRangkingHandler(serviceRank)
	//route
	rank := echo.Group("/rank")
	rank.GET("/", handlerRank.GetUserAndPoint, middleware.JWTMiddleware())
	rank.POST("/create", handlerRank.CreateRangking, middleware.JWTMiddleware())
	rank.GET("/point", handlerRank.GetPointByUserId, middleware.JWTMiddleware())
	rank.POST("/update", handlerRank.UpdateNextUser, middleware.JWTMiddleware())

}
