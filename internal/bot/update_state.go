package bot

import (
	"context"
	"log/slog"
)

func (b *Bot) updateStateFromScreen(ctx context.Context, screen string, filename string) {
	rules := b.Rules[screen]
	newState, err := b.executor.Analyzer().AnalyzeAndUpdateState(b.Gamer, rules, b.Queue)
	if err != nil {
		b.logger.Warn("⚠️ Screen analysis error", slog.String("screen", screen), slog.Any("error", err))
		return
	}

	*b.Gamer = *newState
	b.logger.Info("📥 State updated", slog.String("screen", screen))

	if saveErr := b.Repo.SaveGamer(ctx, newState); saveErr != nil {
		b.logger.Error("❌ Failed to save state.yaml", slog.Any("error", saveErr))
	} else {
		b.logger.Info("💾 Player state saved to state.yaml")
	}
}
