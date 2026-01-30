package main

import (
	"context"
	"log"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/batazor/whiteout-survival-autopilot/internal/bot"
	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/device"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/gift"
	"github.com/batazor/whiteout-survival-autopilot/internal/logger"
	"github.com/batazor/whiteout-survival-autopilot/internal/metrics"
	"github.com/batazor/whiteout-survival-autopilot/internal/redis_queue"
	"github.com/batazor/whiteout-survival-autopilot/internal/repository"
	"github.com/batazor/whiteout-survival-autopilot/internal/syncer"
	"github.com/batazor/whiteout-survival-autopilot/internal/trace"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("❌ Panic caught in main: %v", r)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ─── Initialize OpenTelemetry ──────────────────────────────────────────
	shutdown := trace.Init(ctx, "whiteout-bot")
	defer shutdown()

	// ─── Redis ───────────────────────────────────────────────────────────────
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Redis unavailable: %v", err)
	}

	// ─── Logger ──────────────────────────────────────────────────────────────
	appLogger, err := logger.InitializeLogger("app")
	if err != nil {
		log.Fatalf("❌ Failed to initialize logger: %v", err)
	}

	// ─── Gift listener ───────────────────────────
	gift.AutoStart(gift.Config{
		UserID:      "1634091876319117312",
		DevicesYAML: "db/devices.yaml",
		CodesYAML:   "db/giftCodes.yaml",
		// PythonDir: "",          // script from package
		// PollEvery: 0,           // 0 ⇒ 5 min
		// HistoryDepth: 0,        // 0 ⇒ 10
		Logger: appLogger,
	})

	// ── Metrics ───────────────────────────────────────────────────────────────
	metrics.StartExporter()

	// ─── State repository ─────────────────────────────────────────────────
	repo := repository.NewFileStateRepository("./db/state.yaml")

	// ─── Device / profile configuration ───────────────────────────────────
	devicesCfg, err := config.LoadDeviceConfig("./db/devices.yaml", repo)
	if err != nil {
		log.Fatalf("❌ Configuration loading error: %v", err)
	}

	// 🧠 Update state of all players via Century API
	syncer.RefreshAllPlayersFromCentury(ctx, devicesCfg.AllGamers(), repo, appLogger)

	// ─── Initialize use-cases ─────────────────────────────────────────────
	usecaseLoader := config.NewUseCaseLoader("./usecases")

	// ─── Preload use-cases ────────────────────────────────────────────
	redis_queue.PreloadQueues(ctx, rdb, devicesCfg.AllProfiles(), usecaseLoader)

	// ── Start global task refiller ───────────────────────────────
	go redis_queue.StartGlobalUsecaseRefiller(ctx, devicesCfg, usecaseLoader, rdb, appLogger)

	// ─── Initialize screen analysis rules ───────────────────────────────────────
	rulesUsecases, err := config.LoadAnalyzeRules("references/analyze.yaml")
	if err != nil {
		appLogger.Error("❌ Screen analysis rules loading error", slog.Any("err", err))
		return
	}

	// 🌟 Initialize TriggerEvaluator 🌟
	triggerEvaluator := config.NewTriggerEvaluator()

	// ─── Start devices and bots ────────────────────────────────────────────
	var wg sync.WaitGroup

	for _, devCfg := range devicesCfg.Devices {
		wg.Add(1)

		go func(dc domain.Device) {
			defer wg.Done()

			devLog := appLogger.With("device", dc.Name)

			dev, err := device.New(dc.Name, dc.Profiles, devLog, "./references/area.json", rdb, triggerEvaluator)
			if err != nil {
				devLog.Error("❌ Device creation error", slog.Any("err", err))
				return
			}

			activeGamer, pIdx, gIdx, err := dev.DetectAndSetCurrentGamer(ctx)
			if err != nil || activeGamer == nil {
				devLog.Warn("⚠️ Failed to detect active player", slog.Any("err", err))
				return
			}

			devLog.Info("▶️ Continuing with current player", slog.Int("pIdx", pIdx), slog.Int("gIdx", gIdx), slog.String("nickname", activeGamer.Nickname))

			for {
				select {
				case <-ctx.Done():
					devLog.Info("🛑 Stopping due to context")
					return
				default:
				}

				if pIdx >= len(dc.Profiles) {
					pIdx = 0
				}
				if gIdx >= len(dc.Profiles[pIdx].Gamer) {
					pIdx++
					gIdx = 0
					continue
				}

				target := &dc.Profiles[pIdx].Gamer[gIdx]
				if dev.ActiveGamer() == nil || dev.ActiveGamer().ID != target.ID {
					if err := dev.SwitchTo(ctx, pIdx, gIdx); err != nil {
						devLog.Warn("⚠️ Failed to switch", slog.Any("err", err))
						gIdx++
						continue
					}
				}

				b := bot.NewBot(dev, target, dc.Profiles[pIdx].Email, rdb, rulesUsecases, devLog.With("gamer", target.Nickname), repo)
				b.Play(ctx)

				gIdx++
			}
		}(devCfg)
	}

	wg.Wait()
}
