package redis_queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
)

func PreloadQueues(ctx context.Context, rdb *redis.Client, profiles domain.Profiles, usecaseLoader config.UseCaseLoader) {
	for _, profile := range profiles {
		for _, gamer := range profile.Gamer {
			queue := NewGamerQueue(rdb, gamer.ID)
			key := queue.key()

			// 💣 Complete reset
			if err := rdb.Del(ctx, key).Err(); err != nil {
				fmt.Printf("❌ Failed to clear queue for gamer:%d: %v\n", gamer.ID, err)
				continue
			}

			usecases, err := usecaseLoader.LoadAll(ctx)
			if err != nil {
				fmt.Printf("❌ Error loading usecases for gamer:%d: %v\n", gamer.ID, err)
				continue
			}

			for _, uc := range usecases {
				data, _ := json.Marshal(uc)
				score := float64(100 - uc.Priority)

				if err := rdb.ZAdd(ctx, key, redis.Z{
					Score:  score,
					Member: string(data),
				}).Err(); err != nil {
					fmt.Printf("❌ Failed to add %s to gamer:%d: %v\n", uc.Name, gamer.ID, err)
				} else {
					fmt.Printf("📥 Added usecase %s to gamer:%d\n", uc.Name, gamer.ID)
				}
			}
		}
	}
}
