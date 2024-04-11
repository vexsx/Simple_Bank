package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/vexsx/Simple-Bank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (server *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "session not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find session: %s", err)
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		return nil, status.Errorf(codes.Unauthenticated, "session is blocked: %s", err)
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		return nil, status.Errorf(codes.Unauthenticated, "session is not blonge to user: %s", err)
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		return nil, status.Errorf(codes.Unauthenticated, "missmatch: %s", err)
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		return nil, status.Errorf(codes.Unauthenticated, "session expired: %s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, " %s", err)
	}

	rsp := &pb.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: convertTime(accessPayload.ExpiredAt),
	}
	return rsp, nil
}
