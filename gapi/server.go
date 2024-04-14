package gapi

import (
	"fmt"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/pb"
	"github.com/vexsx/Simple-Bank/token"
	"github.com/vexsx/Simple-Bank/util"
	"github.com/vexsx/Simple-Bank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer serve HTTP
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot  create make token : %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
