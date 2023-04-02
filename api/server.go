package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/piriya-muaithaisong/testgolang_mongo/authorization"
	"github.com/piriya-muaithaisong/testgolang_mongo/db"
	"github.com/piriya-muaithaisong/testgolang_mongo/token"
	"github.com/piriya-muaithaisong/testgolang_mongo/utils"
)

type Server struct {
	config     utils.Config
	store      *db.MongoStore
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store *db.MongoStore) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	authorization.SetupCasbin(config.DBSource)

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("currency", validCurrency)
	// }

	server.setupRouter()

	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	//router.POST("/token/renew_access", server.renewAccessToken)

	// //require auth
	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// authRoutes.POST("/accounts", server.createAccount)
	// authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.GET("/accounts", server.listAccount)

	// authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}
