package data

import (
	"context"
	"math/rand"
	"time"

	grpcclients "github.com/andreamper220/snakeai/internal/server/infrastructure/grpc"
	pb "github.com/andreamper220/snakeai/proto"
)

// Food contains an 'apple' coordinates (starting from 1!).
type Food struct {
	Position Point
}

// CreateRandomFood creates a food object in an empty field cell.
func CreateRandomFood(width, height int, mapID string) *Food {
	x, y := rand.Intn(width), rand.Intn(height)
	if mapID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		gameMap, err := grpcclients.EditorClient.GetMap(ctx, &pb.GetMapRequest{
			Id: mapID,
		})
		if err == nil {
			x, y = checkObstacles(gameMap.GetMap().GetStruct().GetObstacles(), x, y, width, height)
		}
	}

	return &Food{
		Position: Point{X: x + 1, Y: y + 1},
	}
}

func checkObstacles(obstacles []*pb.Obstacle, cx, cy, width, height int) (int, int) {
	for {
		x, y := rand.Intn(width), rand.Intn(height)
		occupied := false
		for _, obstacle := range obstacles {
			if int(obstacle.GetCx()) == x && int(obstacle.GetCy()) == y {
				occupied = true
				break
			}
		}
		if !occupied {
			return x, y
		}
	}
}
