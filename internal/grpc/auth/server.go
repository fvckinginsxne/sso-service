package auth

import (
	"context"
	"errors"

	ssov1 "github.com/fvckinginsxne/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"auth/internal/service/auth"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (*emptypb.Empty, error)
	ParseUserIDFromToken(token string) (int64, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*emptypb.Empty, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	_, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return nil, nil
}

func (s *serverAPI) ParseToken(
	ctx context.Context,
	req *ssov1.ParseTokenRequest,
) (*ssov1.ParseTokenResponse, error) {
	uid, err := s.auth.ParseUserIDFromToken(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid token")
	}

	return &ssov1.ParseTokenResponse{Uid: uid}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
