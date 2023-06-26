package server

import (
	administratorV1 "github.com/ZQCard/kbk-administrator/api/administrator/v1"
	"github.com/ZQCard/kbk-administrator/internal/conf"
	"github.com/ZQCard/kbk-administrator/internal/service"
	"github.com/ZQCard/kbk-administrator/pkg/middleware/requestInfo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *conf.Server, service *service.AdministratorService, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			validate.Validator(),
			recovery.Recovery(),
			tracing.Server(),
			// 元信息
			metadata.Server(),
			// 访问日志
			logging.Server(logger),
			// 租户信息
			requestInfo.SetRequestInfo(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	administratorV1.RegisterAdministratorServiceServer(srv, service)
	return srv
}
