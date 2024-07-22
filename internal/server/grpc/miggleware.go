package internalgrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func (s *Server) LogGRPCRequest(info *grpc.UnaryServerInfo, statusString string, d time.Duration) {
	t := time.Now().Format("02/Jan/2006:15:04:05 -0700")
	msg := fmt.Sprintf("%s %s %s %d", t, info.FullMethod, statusString, d.Microseconds())
	s.logger.Info(msg, "source", "grpc")
}

func (s *Server) RequestLogInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(start)
	s.LogGRPCRequest(info, status.Code(err).String(), duration)
	return result, err
}
