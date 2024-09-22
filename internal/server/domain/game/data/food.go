package data

import (
	"context"
	"math/rand"

	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	pb "github.com/andreamper220/snakeai/proto"
)

type Food struct {
	Position Point
}

func NewFood(width, height int, mapId string) *Food {
	x, y := rand.Intn(width), rand.Intn(height)
	if mapId != "" {
		gameMap, err := grpcclients.EditorClient.GetMap(context.Background(), &pb.GetMapRequest{
			Id: mapId,
		})
		if err == nil {
			x, y = checkObstacles(gameMap.Map.Struct.Obstacles, x, y, width, height)
		}
	}

	return &Food{
		Position: Point{X: x + 1, Y: y + 1},
	}
}

func checkObstacles(obstacles []*pb.Obstacle, cx, cy, width, height int) (int, int) {
	for _, obstacle := range obstacles {
		if int(obstacle.Cx) == cx && int(obstacle.Cy) == cy {
			x, y := rand.Intn(width), rand.Intn(height)
			return checkObstacles(obstacles, x, y, width, height)
		}
	}
	return cx, cy
}
