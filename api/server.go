package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/token"
	"github.com/vexsx/Simple-Bank/util"
)

type Server struct {
	config     util.Config
	store      *db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer serve HTTP
func NewServer(config util.Config, store *db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot  create make token : %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {

	router := gin.Default()
	//user action
	router.POST("/User/Create", server.createUser)
	router.POST("/User/Login", server.loginUser)
	router.GET("/User/:username", server.getUser)

	//from here need auth
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))

	//account actions
	authRoutes.POST("/CreateAccount", server.createAccount)
	authRoutes.GET("/Accounts/:id", server.getAccount)
	authRoutes.GET("/Accounts", server.listAccount)
	authRoutes.POST("/UpdateAccountBalance", server.updateAccountBalance)
	authRoutes.POST("/DeleteAccount", server.deleteAccount)

	//transfers actions
	authRoutes.POST("/Transfer", server.createTransfer)

	server.router = router

}

//start http server

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
