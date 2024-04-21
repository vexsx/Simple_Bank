package gapi

import (
	"context"
	"github.com/lib/pq"
	db "github.com/vexsx/Simple-Bank/db/sqlc"
	"github.com/vexsx/Simple-Bank/pb"
	"github.com/vexsx/Simple-Bank/util"
	"github.com/vexsx/Simple-Bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {

	validation := validateCreateAccountRequest(req)
	if validation != nil {
		return nil, invalidArgumentError(validation)
	}

	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "accounts_owner_fkey", "unique_violation":
				return nil, status.Errorf(codes.PermissionDenied, "%s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	rsp := &pb.CreateAccountResponse{
		Account: convertAccount(account),
	}

	return rsp, nil
}

func validateCreateAccountRequest(req *pb.CreateAccountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	err := val.ValidateCurrency(req.Currency)
	if err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}
	return violations
}
