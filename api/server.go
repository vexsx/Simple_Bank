package api

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/token"
	"github.com/vexsx/Simple-Bank/util"
	"time"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer serve HTTP
func NewServer(config util.Config, store db.Store) (*Server, error) {
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
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}
	server.setUpRouter()
	return server, nil
}

func (server *Server) setUpRouter() {

	router := gin.Default()

	config := cors.Config{
		AllowOrigins:        []string{"*", "http://localhost:4200"},
		AllowMethods:        []string{"PATCH", "POST", "GET"},
		AllowHeaders:        []string{"*", "Origin", "Authorization", "Content-Type", "*"},
		ExposeHeaders:       []string{"Content-Type"},
		AllowPrivateNetwork: true,
		AllowCredentials:    true,
		MaxAge:              12 * time.Hour,
	}

	router.Use(cors.New(config))

	//user action
	router.POST("/User/Create", server.createUser)
	router.POST("/User/Login", server.loginUser)
	router.POST("/Tokens/Renew_Access", server.renewAccessToken)
	router.GET("/User/:username", server.getUser)

	//from here need auth
	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))
	authRoutes.Use(cors.New(config))

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
