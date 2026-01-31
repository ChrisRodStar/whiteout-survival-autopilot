package executor

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/batazor/whiteout-survival-autopilot/internal/adb"
	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
	"github.com/batazor/whiteout-survival-autopilot/internal/metrics"
	"github.com/batazor/whiteout-survival-autopilot/internal/redis_queue"
	"github.com/batazor/whiteout-survival-autopilot/internal/utils"
)

// UseCaseExecutor describes the interface for executing UseCase
type UseCaseExecutor interface {
	ExecuteUseCase(ctx context.Context, uc *domain.UseCase, state *domain.Gamer, queue *redis_queue.Queue)
	Analyzer() Analyzer
}

// Analyzer describes the interface for analyzing screenshots and updating player state
type Analyzer interface {
	AnalyzeAndUpdateState(state *domain.Gamer, rules []domain.AnalyzeRule, queue *redis_queue.Queue) (*domain.Gamer, error)
}

// NewUseCaseExecutor returns an implementation of UseCaseExecutor
func NewUseCaseExecutor(
	logger *slog.Logger,
	triggerEvaluator config.TriggerEvaluator,
	analyzer Analyzer,
	adb adb.DeviceController,
	area *config.AreaLookup,
	botName string,
	queue *redis_queue.Queue,
) UseCaseExecutor {
	return &executorImpl{
		logger:           logger,
		triggerEvaluator: triggerEvaluator,
		analyzer:         analyzer,
		adb:              adb,
		area:             area,
		botName:          botName,
		queue:            queue,
		usecaseLoader:    config.NewUseCaseLoader("./usecases"),
	}
}

type executorImpl struct {
	logger           *slog.Logger
	triggerEvaluator config.TriggerEvaluator
	analyzer         Analyzer
	adb              adb.DeviceController
	area             *config.AreaLookup
	botName          string
	queue            *redis_queue.Queue
	usecaseLoader    config.UseCaseLoader
}

func (e *executorImpl) Analyzer() Analyzer {
	return e.analyzer
}

// ExecuteUseCase executes the entire UseCase
func (e *executorImpl) ExecuteUseCase(ctx context.Context, uc *domain.UseCase, gamer *domain.Gamer, queue *redis_queue.Queue) {
	// Create a span for the entire UseCase
	start := time.Now()
	tracer := otel.Tracer("bot")
	ctx, span := tracer.Start(ctx, uc.Name)
	defer span.End()

	// Extract TraceID for logs
	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()

	// Check UseCase trigger
	if uc.Trigger != "" {
		ok, err := e.triggerEvaluator.EvaluateTrigger(uc.Trigger, gamer)
		if err != nil {
			e.logger.Error("Trigger evaluation failed",
				slog.String("usecase", uc.Name),
				slog.String("trigger", uc.Trigger),
				slog.Any("error", err),
			)
			return
		}

		if !ok {
			e.logger.Warn("Trigger not met, skipping usecase",
				slog.String("usecase", uc.Name),
				slog.String("trigger", uc.Trigger),
			)
			return
		}
	}

	// Логируем старт UseCase с TraceID
	e.logger.Info("=== Start usecase ===",
		slog.String("name", uc.Name),
		slog.String("trace_id", traceID),
	)

	for _, step := range uc.Steps {
		// Call nested steps
		e.runStep(ctx, step, 0, gamer)
	}

	// Log UseCase completion with TraceID
	e.logger.Info("=== End usecase ===",
		slog.String("name", uc.Name),
		slog.String("trace_id", traceID),
	)

	// If UseCase executed successfully — set TTL (if exists)
	if uc.TTL > 0 && queue != nil {
		if err := queue.SetLastExecuted(ctx, gamer.ID, uc.Name, uc.TTL); err != nil {
			e.logger.Error("Failed to set last executed TTL", slog.Any("error", err))
		}
	}

	// Counters and metrics
	metrics.UsecaseTotal.WithLabelValues(uc.Name).Inc()
	metrics.UsecaseDuration.WithLabelValues(uc.Name).Observe(time.Since(start).Seconds())

	// Example of recording player state metrics
	if gamer != nil {
		// Player power
		metrics.GamerPowerGauge.WithLabelValues(gamer.Nickname).Set(float64(gamer.Power))

		// Furnace level (if available)
		if gamer.Buildings.Furnace.Level > 0 {
			metrics.GamerFurnaceLevel.WithLabelValues(gamer.Nickname).Set(float64(gamer.Buildings.Furnace.Level))
		}
	}
}

// runStep executes one UseCase step (possibly recursively calling itself for nested steps)
func (e *executorImpl) runStep(ctx context.Context, step domain.Step, indent int, gamer *domain.Gamer) bool {
	// Start tracing for each step
	ctx, stepSpan := otel.Tracer("bot").Start(ctx, "runStep: "+step.Action)
	defer stepSpan.End()

	select {
	case <-ctx.Done():
		e.logger.Warn("Step cancelled by context")
		return true
	default:
	}

	prefix := strings.Repeat("  ", indent)

	// If step.Click exists — click
	if step.Click != "" {
		e.logger.Info(prefix+"Click", slog.String("target", step.Click))

		err := e.adb.ClickRegion(step.Click, e.area)
		if err != nil {
			e.logger.Error(prefix+"Failed to click region",
				slog.String("target", step.Click),
				slog.Any("error", err),
			)
			return true
		}
	}

	// If step.Action exists — execute it
	if step.Action != "" {
		e.logger.Info(prefix+"Click", slog.String("action", step.Action))

		switch step.Action {

		// Reset state field: "reset"
		case "reset":
			if step.Set == "" {
				e.logger.Warn(prefix + "Reset skipped: missing 'set' field")
				return false
			}

			// Get current value for logging
			prevVal, _ := utils.GetStateFieldByPath(gamer, step.Set)

			if err := utils.SetStateFieldByPath(gamer, step.Set, step.To); err != nil {
				e.logger.Error(prefix+"Failed to reset state field",
					slog.String("path", step.Set),
					slog.Any("from", prevVal),
					slog.Any("to", step.To),
					slog.Any("error", err),
				)
			} else {
				e.logger.Info(prefix+"State field reset",
					slog.String("path", step.Set),
					slog.Any("from", prevVal),
					slog.Any("to", step.To),
				)
			}

		// Loop organization: "loop"
		case "loop":
			if step.Trigger == "" {
				e.logger.Warn(prefix + "Loop trigger is missing, skipping loop")
				return false
			}

			// Create a separate span for the entire loop
			loopCtx, loopSpan := otel.Tracer("bot").Start(ctx, prefix+"loop: "+step.Trigger)
			defer loopSpan.End()

			e.logger.Info(prefix+"Entering loop", slog.String("trigger", step.Trigger))

			for {
				select {
				case <-loopCtx.Done():
					e.logger.Warn(prefix + "Loop interrupted by context")
					return true
				default:
				}

				shouldContinue, err := e.triggerEvaluator.EvaluateTrigger(step.Trigger, gamer)
				if err != nil {
					e.logger.Error(prefix+"Trigger evaluation failed", slog.Any("error", err))
					break
				}
				if !shouldContinue {
					e.logger.Info(prefix + "Loop trigger returned false, exiting loop")
					break
				}

				for _, s := range step.Steps {
					if stopped := e.runStep(loopCtx, s, indent+1, gamer); stopped {
						e.logger.Info(prefix + "Loop stopped manually (loop_stop)")
						return false
					}
				}
			}

		// Forced loop exit
		case "loop_stop":
			e.logger.Info(prefix + "Received loop_stop")
			return true

		// Screenshot with subsequent analysis
		case "screenshot":
			// If there are analysis rules
			if len(step.Analyze) > 0 {
				_, analyzeSpan := otel.Tracer("bot").Start(ctx, prefix+"AnalyzeAndUpdateState")
				defer analyzeSpan.End()

				newState, err := e.analyzer.AnalyzeAndUpdateState(gamer, step.Analyze, e.queue)
				if err != nil {
					e.logger.Error(prefix+"Analyze failed", slog.Any("error", err))
				} else {
					*gamer = *newState
					e.logger.Info(prefix + "Analyze completed and state updated")
				}
			}
		}
	}

	// If step.Wait exists — wait
	if step.Wait > 0 {
		e.logger.Info(prefix+"Wait", slog.Duration("duration", step.Wait))
		select {
		case <-time.After(step.Wait):
		case <-ctx.Done():
			e.logger.Warn(prefix+"Wait interrupted by context cancel", slog.Duration("wait", step.Wait))
			return true
		}
	}

	// If there's an if/then/else condition
	if step.If != nil {
		// Create a separate span for the entire `if`
		ifCtx, ifSpan := otel.Tracer("bot").Start(ctx, prefix+"if: "+step.If.Trigger)
		defer ifSpan.End()

		e.logger.Info(prefix+"If Trigger", slog.String("expr", step.If.Trigger))

		result, err := e.triggerEvaluator.EvaluateTrigger(step.If.Trigger, gamer)
		if err != nil {
			e.logger.Error(prefix+"Trigger evaluation failed",
				slog.String("expr", step.If.Trigger),
				slog.Any("error", err),
			)
			return false
		}

		if result {
			// then
			thenCtx, thenSpan := otel.Tracer("bot").Start(ifCtx, prefix+"then")
			defer thenSpan.End()

			e.logger.Info(prefix + "Condition met: executing THEN")
			for _, s := range step.If.Then {
				stopped := e.runStep(thenCtx, s, indent+1, gamer)
				if stopped {
					return true
				}
			}
		} else if len(step.If.Else) > 0 {
			// else
			elseCtx, elseSpan := otel.Tracer("bot").Start(ifCtx, prefix+"else")
			defer elseSpan.End()

			e.logger.Info(prefix + "Condition NOT met: executing ELSE")
			for _, s := range step.If.Else {
				stopped := e.runStep(elseCtx, s, indent+1, gamer)
				if stopped {
					return true
				}
			}
		}
	}

	// Long tap (longtap)
	if step.Longtap != "" {
		e.logger.Info(prefix+"Longtap", slog.String("target", step.Longtap), slog.Duration("hold", step.Wait))

		bbox, err := e.area.GetRegionByName(step.Longtap)
		if err != nil {
			e.logger.Error(prefix+"Failed to find region for longtap",
				slog.String("target", step.Longtap),
				slog.Any("error", err),
			)
			return true
		}

		x, y, _, _ := bbox.ToPixels()
		err = e.adb.Swipe(x, y, x, y, step.Wait) // swipe to the same place with specified time
		if err != nil {
			e.logger.Error(prefix+"Failed to perform longtap",
				slog.String("target", step.Longtap),
				slog.Any("error", err),
			)
			return true
		}
	}

	// --- PUSH-USECASE --------------------------------------------
	if len(step.PushUsecase) > 0 && e.queue != nil {
		for _, push := range step.PushUsecase {
			// 1) check trigger (if exists)
			if push.Trigger != "" {
				ok, err := e.triggerEvaluator.EvaluateTrigger(push.Trigger, gamer)
				if err != nil {
					e.logger.Error("Trigger evaluation failed for pushUsecase",
						slog.String("trigger", push.Trigger), slog.Any("error", err))
					continue
				}
				if !ok {
					e.logger.Debug("pushUsecase trigger not satisfied",
						slog.String("trigger", push.Trigger))
					continue
				}
			}

			// If trigger is satisfied, add usecase to queue
			for _, uc := range push.List {
				ucOriginal := e.usecaseLoader.GetByName(uc.Name)

				e.logger.Info("📥 Push usecase from analysis", slog.String("usecase", uc.Name))
				if err := e.queue.Push(context.Background(), ucOriginal); err != nil {
					e.logger.Error("❌ Failed to push usecase", slog.String("usecase", uc.Name), slog.Any("error", err))
				}
			}
		}
	}

	return false
}
