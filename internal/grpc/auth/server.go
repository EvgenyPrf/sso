package auth

import (
	"context"
	ssov1 "github.com/EvgenyPrf/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	//заглушка, если не хочется реализовывать методы интерфейса
	ssov1.UnimplementedAuthServer
}

// регистрация хендлеров
func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me")
}
func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	panic("implement me")
}
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me")
}
