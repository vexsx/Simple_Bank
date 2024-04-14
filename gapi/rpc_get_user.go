package gapi

import (
	"context"
	"database/sql"
	"github.com/vexsx/Simple-Bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.Internal, "user dosent exist :  %s ", err)
		}
		return nil, status.Errorf(codes.Internal, "cannot find user :  %s ", err)
	}

	rsp := &pb.GetUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
