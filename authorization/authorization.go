package authorization

import (
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
)

//SetupRoutes : all the routes are defined here
func SetupCasbin(mongoUrl string) {

	// Initialize  casbin adapter
	adapter, err := mongodbadapter.NewAdapter(mongoUrl)
	if err != nil {
		panic(err)
	}
	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("./authorization/config/rbac_model.conf", adapter)
	if err != nil {
		panic(err)
	}

	//add policy
	if hasPolicy := enforcer.HasPolicy("doctor", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("doctor", "report", "read")
	}
	if hasPolicy := enforcer.HasPolicy("doctor", "report", "write"); !hasPolicy {
		enforcer.AddPolicy("doctor", "report", "write")
	}
	if hasPolicy := enforcer.HasPolicy("patient", "report", "read"); !hasPolicy {
		enforcer.AddPolicy("patient", "report", "read")
	}

	// userRepository := repository.NewUserRepository(db)

	// if err := userRepository.Migrate(); err != nil {
	// 	log.Fatal("User migrate err", err)
	// }

	// userController := controller.NewUserController(userRepository)

	// apiRoutes := httpRouter.Group("/api")

	// {
	// 	apiRoutes.POST("/register", userController.AddUser(enforcer))
	// 	apiRoutes.POST("/signin", userController.SignInUser)
	// }

	// userProtectedRoutes := apiRoutes.Group("/users", middleware.AuthorizeJWT())
	// {
	// 	userProtectedRoutes.GET("/", middleware.Authorize("report", "read", enforcer), userController.GetAllUser)
	// 	userProtectedRoutes.GET("/:user", middleware.Authorize("report", "read", enforcer), userController.GetUser)
	// 	userProtectedRoutes.PUT("/:user", middleware.Authorize("report", "write", enforcer), userController.UpdateUser)
	// 	userProtectedRoutes.DELETE("/:user", middleware.Authorize("report", "write", enforcer), userController.DeleteUser)
	// }

	// httpRouter.Run()

}
