package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer serve HTTP
func NewServer(store *db.Store) *Server {

	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil
		}
	}
	//user action
	router.POST("/CreateUser", server.createUser)
	router.GET("/User/:username", server.getUser)

	//account actions
	router.POST("/CreateAccount", server.createAccount)
	router.GET("/Accounts/:id", server.getAccount)
	router.GET("/Accounts", server.listAccount)
	router.POST("/UpdateAccountBalance", server.updateAccountBalance)
	router.POST("/DeleteAccount", server.deleteAccount)

	//transfers actions
	router.POST("/Transfer", server.createTransfer)

	server.router = router
	return server
}

//start http server

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
