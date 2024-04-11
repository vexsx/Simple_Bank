package gapi

import (
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}

func convertTime(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
