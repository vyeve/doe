package domainserver

import (
	"context"
	"io"
	"net"
	"os"

	"doe/source/data"
	"doe/source/logger"
	"doe/source/models"

	gRPC "google.golang.org/grpc"
)

type Server interface {
	Serve() error
}

func New(params Params) Server {
	srv := &serverImpl{
		log: params.Logger,
	}
	grpcServer := gRPC.NewServer()
	models.RegisterPortServiceServer(grpcServer, &serverClient{
		repo: params.Repo,
		log:  params.Logger,
	})
	srv.srv = grpcServer
	srv.log.Infof("Domain server. Start listening on :%s", "localhost:9098")
	return srv
}

type serverImpl struct {
	log logger.Logger
	srv *gRPC.Server
}

func (s *serverImpl) Serve() error {
	addr := os.Getenv(HostURLEnvKey)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.log.Errorf("failed to listen: %v", err)
		return err
	}
	s.log.Infof("gRPC server started to listen at: %s", addr)
	return s.srv.Serve(lis)
}

type serverClient struct {
	models.UnimplementedPortServiceServer
	repo data.Repository
	log  logger.Logger
}

// GetAll returns All data from storage to stream
func (s serverClient) GetAll(ctx context.Context, r *models.Request) (*models.Ports, error) {
	ports, err := s.repo.GetAll(context.Background(), r.Limit)
	if err != nil {
		s.log.Warnf("Failed to extract ports")
		return nil, err
	}
	result := new(models.Ports)
	result.Ports = ports
	return result, nil
}

// GetOne returns one port information by port ID
func (s serverClient) GetOne(ctx context.Context, r *models.Request) (*models.Port, error) {
	// defer s.dbClient.CloseSession()
	return s.repo.GetByID(ctx, r.PortID)
}

// Set uploads data from stream to storage
func (s serverClient) Set(stream models.PortService_SetServer) error {
	ports := make([]*models.Port, 0, insertLimit)
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		ports = append(ports, data)
		if len(ports) == insertLimit {
			err = s.repo.Insert(context.Background(), ports)
			if err != nil {
				return err
			}
			ports = ports[0:0]
		}
	}
	if len(ports) > 0 {
		err := s.repo.Insert(context.Background(), ports)
		if err != nil {
			return err
		}
	}
	return stream.SendAndClose(nil)
}
