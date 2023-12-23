package authgrpc

import (
	"context"
	"errors"

	"github.com/1kovalevskiy/sso/internal/entity"
	error_ "github.com/1kovalevskiy/sso/internal/error"

	ssov1 "github.com/1kovalevskiy/proto_sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	GetCreateApp(ctx context.Context, name string, password string, secret string, ttlHour int) (int, error)
	Login(ctx context.Context, email string, password string, appID int) (string, error)
	RegisterNewUser(ctx context.Context, email string, pass string) (int, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func New(auth Auth) func(gRPCServer *grpc.Server) {
	return func(gRPCServer *grpc.Server) {
		ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
	}
}

func (s *serverAPI) AddApp(ctx context.Context, in *ssov1.AddAppRequest) (*ssov1.AddAppResponse, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "secret is required")
	}

	if in.GetTtlHour() == 0 {
		return nil, status.Error(codes.InvalidArgument, "ttl_hour is required")
	}

	id, err := s.auth.GetCreateApp(ctx, in.GetName(), in.GetPassword(), in.GetSecret(), int(in.GetTtlHour()))
	if err != nil {
		if errors.Is(err, error_.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid name or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.AddAppResponse{AppId: int32(id)}, nil
}

func (s *serverAPI) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword(), int(in.GetAppId()))
	if err != nil {
		if errors.Is(err, error_.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	uid, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{UserId: int64(uid)}, nil
}
