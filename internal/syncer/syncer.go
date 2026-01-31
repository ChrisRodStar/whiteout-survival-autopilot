package syncer

import (
	"context"
	"log/slog"
	"sync"

	"github.com/batazor/whiteout-survival-autopilot/internal/century"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/repository"
)

// RefreshAllPlayersFromCentury loads data for all players via Century API and saves them to state.yaml
func RefreshAllPlayersFromCentury(
	ctx context.Context,
	gamers []*domain.Gamer,
	repo repository.StateRepository,
	logger *slog.Logger,
) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var updatedGamers []domain.Gamer

	for _, g := range gamers {
		gamer := g
		wg.Add(1)

		go func() {
			defer wg.Done()

			info, err := century.FetchPlayerInfo(gamer.ID)
			if err != nil {
				logger.Warn("⚠️ Failed to get player data from Century", slog.Int("id", gamer.ID), slog.Any("err", err))
				return
			}

			// Update data
			gamer.Nickname = info.Data.Nickname
			gamer.State = info.Data.KID
			gamer.Avatar = info.Data.AvatarImage
			gamer.Buildings.Furnace.Level = info.Data.StoveLevel

			mu.Lock()
			updatedGamers = append(updatedGamers, *gamer)
			mu.Unlock()

			logger.Info("📥 Player updated from Century", slog.String("nickname", gamer.Nickname), slog.Int("id", gamer.ID))
		}()
	}

	wg.Wait()

	// 💾 Save final state.yaml
	finalState := &domain.State{Gamers: updatedGamers}
	if err := repo.SaveState(ctx, finalState); err != nil {
		logger.Error("❌ Failed to save state.yaml after update", slog.Any("error", err))
	} else {
		logger.Info("💾 Final state.yaml successfully saved")
	}
}
