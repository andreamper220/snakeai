package application

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"strconv"

	"github.com/andreamper220/snakeai/internal/editor/domain"
	"github.com/andreamper220/snakeai/internal/editor/infrastructure/storages"
	"github.com/andreamper220/snakeai/pkg/logger"
	pb "github.com/andreamper220/snakeai/proto"
)

const (
	MinWidth  = 5
	MaxWidth  = 30
	MinHeight = 5
	MaxHeight = 30
)

var srv *grpc.Server

type EditorServer struct {
	pb.UnimplementedEditorServer
}

func (s EditorServer) CheckMap(ctx context.Context, request *pb.CheckMapRequest) (*pb.CheckMapResponse, error) {
	width := request.Struct.Width
	height := request.Struct.Height
	if width < MinWidth || width > MaxWidth {
		return nil, status.Error(codes.InvalidArgument, "width out of range")
	}
	if height < MinHeight || (height > MaxHeight) {
		return nil, status.Error(codes.InvalidArgument, "height out of range")
	}
	return &pb.CheckMapResponse{}, nil
}
func (s EditorServer) SaveMap(ctx context.Context, request *pb.SaveMapRequest) (*pb.SaveMapResponse, error) {
	requestObstacles := request.Struct.GetObstacles()
	obstacles := make([][2]int32, len(requestObstacles))
	for i := 0; i < len(requestObstacles); i++ {
		obstacles[i][0] = requestObstacles[i].Cx
		obstacles[i][1] = requestObstacles[i].Cy
	}
	gameMap := &domain.Map{
		Width:     request.Struct.Width,
		Height:    request.Struct.Height,
		Obstacles: obstacles,
	}

	mapId, err := storages.EditorStorage.AddMap(gameMap)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.SaveMapResponse{
		Map: &pb.Map{
			Id:     mapId.String(),
			Struct: request.Struct,
		},
	}, nil
}
func (s EditorServer) GetMap(ctx context.Context, request *pb.GetMapRequest) (*pb.GetMapResponse, error) {
	mapId, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	gameMap, err := storages.EditorStorage.GetMap(mapId)
	if err != nil {
		if errors.Is(err, storages.ErrMapNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	obstacles := make([]*pb.Obstacle, len(gameMap.Obstacles))
	for i := 0; i < len(obstacles); i++ {
		obstacle := &pb.Obstacle{
			Cx: gameMap.Obstacles[i][0],
			Cy: gameMap.Obstacles[i][1],
		}
		obstacles[i] = obstacle
	}

	responseMap := &pb.Map{
		Id: gameMap.Id.String(),
		Struct: &pb.MapStruct{
			Width:     gameMap.Width,
			Height:    gameMap.Height,
			Obstacles: obstacles,
		},
	}
	return &pb.GetMapResponse{
		Map: responseMap,
	}, nil
}

func MakeStorage() error {
	storages.EditorStorage = storages.NewMemStorage()

	return nil
}

func Run(port int, serverless bool) error {
	if err := logger.Initialize(); err != nil {
		return err
	}
	logger.Log.Info("logger established")

	if err := MakeStorage(); err != nil {
		return err
	}

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		logger.Log.Fatal(err)
	}

	if !serverless {
		srv = InitGRPCServer(port)
	}
	logger.Log.Infof("gRPC server listening on port %d", port)
	if err = srv.Serve(listen); err != nil {
		logger.Log.Fatal("gRPC server Serve: %v", err)
		return err
	}
	return nil
}

func Stop() {
	srv.Stop()
}

func InitGRPCServer(port int) *grpc.Server {
	srv = grpc.NewServer()
	pb.RegisterEditorServer(srv, &EditorServer{})

	return srv
}
