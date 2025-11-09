package auth

import (
	"context"

	"buf.build/go/protovalidate"
	ssov1 "github.com/EvgenyPrf/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email,
		password string,
		appID int,
	) (token string, err error)

	RegisterNewUser(ctx context.Context,
		email,
		password string,
	) (userId int64, err error)

	IsAdmin(ctx context.Context, userId int64) (bool, error)
}
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	v    protovalidate.Validator
	auth Auth
}

// Register регистрация хендлеров и инициализация валидатора
func Register(gRPC *grpc.Server, auth Auth) {
	v, err := protovalidate.New()
	if err != nil {
		// В проде лучше вернуть ошибку наружу, а не паниковать
		panic("protovalidate init: " + err.Error())
	}
	ssov1.RegisterAuthServer(gRPC, &serverAPI{v: v, auth: auth})
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := s.v.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid RegisterRequest: %v", err)
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &ssov1.RegisterResponse{UserId: userId}, nil
}

// Login — валидация через Protovalidate + простая логика
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := s.v.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid LoginRequest: %v", err)
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := s.v.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid IsAdminRequest: %v", err)
	}

	res, err := s.auth.IsAdmin(ctx, req.GetUserId())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: res}, nil
}
