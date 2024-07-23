//go:build integration
// +build integration

package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/OlgaResh1/OtusGoHomeWork/SystemMonitoring/internal/pb"
	suite "github.com/stretchr/testify/suite"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/durationpb"
)

type StatSuite struct {
	suite.Suite
	ctx         context.Context
	serviceConn *grpc.ClientConn
	statClient  pb.ServiceStatClient
}

func (s *StatSuite) SetupSuite() {
	servceAddress := os.Getenv("SYSTEM_MONITORING_ADDR")
	if servceAddress == "" {
		servceAddress = "127.0.0.1:50051"
	}

	conn, err := grpc.NewClient(servceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().Nil(err)

	s.ctx = context.Background()
	s.serviceConn = conn
	s.statClient = pb.NewServiceStatClient(s.serviceConn)
}

func (s *StatSuite) SetupTest() {
	// Setup test
}

func (s *StatSuite) getRequest() (*pb.ServiceStatRequest, error) {
	interval, err := time.ParseDuration("5s")
	if err != nil {
		return nil, err
	}
	avgWindow, err := time.ParseDuration("15s")
	if err != nil {
		return nil, err
	}
	return &pb.ServiceStatRequest{
		Interval:  durationpb.New(interval),
		Avgwindow: durationpb.New(avgWindow),
	}, nil
}

func (s *StatSuite) getBadRequest() (*pb.ServiceStatRequest, error) {
	interval, err := time.ParseDuration("10ms")
	if err != nil {
		return nil, err
	}
	avgWindow, err := time.ParseDuration("1s")
	if err != nil {
		return nil, err
	}
	return &pb.ServiceStatRequest{
		Interval:  durationpb.New(interval),
		Avgwindow: durationpb.New(avgWindow),
	}, nil
}

func (s *StatSuite) TestReceiveStat() {
	req, err := s.getRequest()
	s.Require().NoError(err)
	stream, err := s.statClient.GetFullStatStream(s.ctx, req)
	s.Require().NoError(err)
	stat, err := stream.Recv()
	s.Require().NoError(err)
	s.Require().NotNil(stat)
}

func (s *StatSuite) TestErrorInterval() {
	req, err := s.getBadRequest()
	s.Require().NoError(err)
	stream, err := s.statClient.GetFullStatStream(s.ctx, req)
	s.Require().NoError(err)
	stat, err := stream.Recv()
	s.Require().Error(err)
	s.Require().Nil(stat)
}

func (s *StatSuite) TearDownTest() {
	// Teardown test
}

func (s *StatSuite) TearDownSuite() {
	err := s.serviceConn.Close()
	if err != nil {
		s.T().Log(err)
	}
}

func TestStatSuite(t *testing.T) {
	suite.Run(t, new(StatSuite))
}
