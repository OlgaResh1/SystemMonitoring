//go:generate protoc --proto_path=../../../api/ --go_out=../../../internal/pb --go-grpc_out=../../../internal/pb ../../../api/SystemStats.proto

package internalgrpc

import (
	"fmt"
	"net"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/config"
	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedServiceStatServer
	logger      Logger
	app         Application
	server      *grpc.Server
	addr        string
	minInterval time.Duration
	maxWindow   time.Duration
}

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
}

type Application interface {
	StatByIntervalPB(queueType int, avgWindow time.Duration) (any, error)
	FullStatByIntervalPB(avgWindow time.Duration) (any, error)
}

func NewServer(cfg config.Config, logger Logger, app Application) *Server {
	return &Server{
		logger:      logger,
		app:         app,
		addr:        cfg.GRPC.Address,
		minInterval: cfg.GRPC.MinInterval,
		maxWindow:   cfg.Stat.MaxAvgWindow,
	}
}

func (s *Server) Start() error {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			s.RequestLogInterceptor,
		),
	)
	s.server = server

	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.logger.Info(fmt.Sprintf("GRPC Server started %s", s.addr))
	pb.RegisterServiceStatServer(s.server, s)
	return s.server.Serve(lsn)
}

func (s *Server) Stop() error {
	if s.server != nil {
		s.server.GracefulStop()
	}
	return nil
}

func (s *Server) GetFullStatStream(req *pb.ServiceStatRequest, srv pb.ServiceStat_GetFullStatStreamServer) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "intervals is not specified")
	}
	interval := req.GetInterval().AsDuration()
	avgWindow := req.GetAvgwindow().AsDuration()

	if interval < s.minInterval || interval <= 0 {
		return status.Error(codes.InvalidArgument, "interval value not valid")
	}
	if avgWindow > s.maxWindow || avgWindow <= 0 {
		return status.Error(codes.InvalidArgument, "avgWindow value not valid")
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-srv.Context().Done():
			s.logger.Info("context done", "source", "grpc")
			return nil

		case <-ticker.C:

			msg := &pb.ServiceStatResponse{
				FullStat: s.GetFullStat(req.GetInterval(), avgWindow),
			}
			if err := srv.Send(msg); err != nil {
				s.logger.Error("unable to send message", "error", err, "source", "grpc")
				return err
			}
		}
	}
}

func (s *Server) GetFullStat(interval *durationpb.Duration, avgWindow time.Duration) *pb.FullStat {
	stat, err := s.app.FullStatByIntervalPB(avgWindow)
	if err != nil {
		s.logger.Error("error getting stats", "error", err, "source", "grpc")
		return nil
	}
	fullStat, ok := stat.(*pb.FullStat)
	if !ok {
		s.logger.Error(fmt.Sprintf("invalid response type: expected pb.FullStat, got %T", stat), "source", "grpc")
		return nil
	}
	fullStat.Time = timestamppb.Now()
	fullStat.Interval = interval
	return fullStat
}
