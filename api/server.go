package api

import (
	"github.com/gin-gonic/gin"
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

	router.POST("/CreateAccount", server.createAccount)
	router.GET("/Accounts/:id", server.getAccount)
	router.GET("/Accounts", server.listAccount)
	router.POST("/UpdateAccountBalance", server.updateAccountBalance)
	router.POST("/DeleteAccount", server.deleteAccount)

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
