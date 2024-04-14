package gapi

import (
	"context"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateAccountBalance(ctx context.Context, req *pb.UpdateAccountBalanceRequest) (*pb.UpdateAccountBalanceResponse, error) {

	_, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.AddAccountBalanceParams{
		ID:     req.GetId(),
		Amount: req.GetAmount(),
	}

	account, err := server.store.AddAccountBalance(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, ": %s", err)
	}

	rsp := &pb.UpdateAccountBalanceResponse{
		Account: convertAccount(account),
	}
	return rsp, nil
}
