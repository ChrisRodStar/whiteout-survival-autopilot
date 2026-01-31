package bot

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/batazor/whiteout-survival-autopilot/internal/analyzer"
	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/device"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain/state"
	"github.com/batazor/whiteout-survival-autopilot/internal/executor"
	"github.com/batazor/whiteout-survival-autopilot/internal/fsm"
	"github.com/batazor/whiteout-survival-autopilot/internal/redis_queue"
	"github.com/batazor/whiteout-survival-autopilot/internal/repository"
)

type Bot struct {
	Gamer    *domain.Gamer
	Email    string
	Device   *device.Device
	Queue    *redis_queue.Queue
	logger   *slog.Logger
	Rules    config.ScreenAnalyzeRules
	Repo     repository.StateRepository
	executor executor.UseCaseExecutor
}

func NewBot(dev *device.Device, gamer *domain.Gamer, email string, rdb *redis.Client, rules config.ScreenAnalyzeRules, log *slog.Logger, repo repository.StateRepository) *Bot {
	queue := redis_queue.NewGamerQueue(rdb, gamer.ID)

	exec := executor.NewUseCaseExecutor(
		log,
		config.NewTriggerEvaluator(),
		analyzer.NewAnalyzer(dev.AreaLookup, log, dev.OCRClient),
		dev.ADB,
		dev.AreaLookup,
		gamer.Nickname,
		queue,
	)

	return &Bot{
		Gamer:    gamer,
		Email:    email,
		Device:   dev,
		Queue:    queue,
		logger:   log,
		Rules:    rules,
		Repo:     repo,
		executor: exec,
	}
}

func (b *Bot) Play(ctx context.Context) {
	// 📸 Analyze state on the main screen
	b.updateStateFromScreen(ctx, "main_city", "out/bot_"+b.Gamer.Nickname+"_start_main_city.png")

	for {
		select {
		case <-ctx.Done():
			b.logger.Warn("🛑 Context cancelled — stopping bot")
			return
		default:
		}

		// get use-case from queue
		uc, err := b.Queue.PopBest(ctx, b.Gamer.ScreenState.CurrentState)
		if err != nil {
			b.logger.Warn("⚠️ Failed to get use-case", "err", err)
			continue
		}

		// queue is empty → exit to switch to another player
		if uc == nil {
			b.logger.Info("📭 Queue is empty — stopping bot")
			break
		}

		// 🕒 Check TTL (skip usecase if not expired)
		shouldSkip, err := b.Queue.ShouldSkip(ctx, b.Gamer.ID, uc.Name)
		if err != nil {
			b.logger.Error("❌ Failed to check usecase TTL", slog.Any("err", err))
			continue
		}
		if shouldSkip {
			b.logger.Info("⏭️ UseCase skipped by TTL", slog.String("name", uc.Name))
			continue
		}

		b.logger.Info("🚀 Executing use-case", "name", uc.Name, "priority", uc.Priority)

		// switch to the usecase start screen
		switchedScreen := false
		if b.Gamer.ScreenState.CurrentState != uc.Node {
			b.logger.Info("🔁 Switching to usecase screen", slog.String("name", uc.Name), slog.String("screen", uc.Node))
			errForceTo := b.Device.FSM.ForceTo(uc.Node, b.updateStateFromScreen)
			if errForceTo != nil {
				if errors.Is(errForceTo, fsm.EventNotActive) {
					b.logger.Info("⏭️ UseCase skipped because event is not active", slog.String("name", uc.Name))

					// Set TTL for usecase in queue
					errSetLastExecuted := b.Queue.SetLastExecuted(ctx, b.Gamer.ID, uc.Name, uc.TTL)
					if errSetLastExecuted != nil {
						b.logger.Error("❌ Failed to set usecase TTL", slog.Any("err", err))
					}

					continue
				}

				b.logger.Error("❌ Failed to switch to usecase screen", slog.Any("err", errForceTo))
			} else {
				switchedScreen = true
			}
		} else {
			b.logger.Info("🔁 Already on usecase screen", slog.String("name", uc.Name), slog.String("screen", uc.Node))
		}

		// Call updateStateFromScreen only if FSM didn't do it in ForceTo, or if there was no transition
		if !switchedScreen {
			b.updateStateFromScreen(ctx, uc.Node, "out/bot_"+b.Gamer.Nickname+"_before_trigger.png")
		}

		b.executor.ExecuteUseCase(ctx, uc, b.Gamer, b.Queue)

		// Time for screen rendering
		time.Sleep(1 * time.Second)
	}

	// Time for screen rendering
	time.Sleep(2 * time.Second)

	// 🔁 Return to main screen
	b.Device.FSM.ForceTo(state.StateMainCity, nil)

	// Time for screen rendering
	time.Sleep(2 * time.Second)

	b.logger.Info("⏭️ Queue completed. Ready to switch.")
}
