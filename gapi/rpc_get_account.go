package gapi

import (
	"context"
	"database/sql"
	"errors"
	"github.com/vexsx/Simple-Bank/pb"
	"github.com/vexsx/Simple-Bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {

	account, err := server.store.GetAccount(ctx, req.GetId())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, ": %s", err)
		}
		return nil, status.Errorf(codes.Internal, ": %s", err)
	}

	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		return nil, status.Errorf(codes.Unauthenticated, ": %s", err)
	}

	rsp := &pb.GetAccountResponse{
		Account: convertAccount(account),
	}
	return rsp, nil
}
