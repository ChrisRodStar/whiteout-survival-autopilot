package redis_queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/batazor/whiteout-survival-autopilot/internal/config"
	"github.com/batazor/whiteout-survival-autopilot/internal/domain"
)

const (
	// screenBoost – how many "virtual" points a UseCase gets whose Node
	// matches the current screen. Higher — more aggressively process tasks
	// of the current screen before moving to others.
	screenBoost = 5

	// scanWindow – how many top ZSET elements we look at before choosing
	// the best. 10–50 is enough to almost always "see" all UCs of the current
	// screen without overloading Redis.
	scanWindow = 20
)

// Queue stores tasks for a specific bot (or gamer) in
// a Redis priority queue (sorted‑set).
//   - The higher uc.Priority (0–100), the lower the score, the further left the element in ZSET.
//   - PopBest additionally shifts the score of "its" screen by screenBoost.
//
// Value format = json.Marshal(domain.UseCase).
// Score is calculated on Push: score = 100 - uc.Priority.
// -----------------------------------------------------------------------------
type Queue struct {
	rdb   *redis.Client
	botID string
}

func (q *Queue) key() string {
	return fmt.Sprintf("bot:queue:%s", q.botID)
}

func NewGamerQueue(rdb *redis.Client, gamerID int) *Queue {
	return &Queue{rdb: rdb, botID: fmt.Sprintf("gamer:%d", gamerID)}
}

// Push adds a UseCase to the Redis priority queue.
// The higher uc.Priority (0–100), the higher the task priority (lower score).
func (q *Queue) Push(ctx context.Context, uc *domain.UseCase) error {
	data, err := json.Marshal(uc)
	if err != nil {
		return err
	}

	score := float64(100 - uc.Priority) // Higher priority, lower score
	return q.rdb.ZAdd(ctx, q.key(), redis.Z{
		Score:  score,
		Member: data,
	}).Err()
}

// Pop extracts the highest priority UseCase from the Redis queue.
func (q *Queue) Pop(ctx context.Context) (*domain.UseCase, error) {
	items, err := q.rdb.ZPopMin(ctx, q.key(), 1).Result()
	if err != nil || len(items) == 0 {
		return nil, err
	}

	var uc domain.UseCase
	if err := json.Unmarshal([]byte(items[0].Member.(string)), &uc); err != nil {
		return nil, err
	}

	return &uc, nil
}

// -----------------------------------------------------------------------------
// PopBest – extracts UC considering current screen "boost".
// -----------------------------------------------------------------------------

// currentNode – FSM node where the UI is currently located (e.g., "alliance").
//
// Algorithm:
//  1. Take the first scanWindow elements from ZSET (they are already roughly sorted
//     by priority).
//  2. For each, decode UseCase, calculate adjustedScore = score - screenBoost,
//     if uc.Node == currentNode.
//  3. Choose the minimum adjustedScore.
//  4. Remove this element from ZSET (ZREM) and return.
func (q *Queue) PopBest(ctx context.Context, currentNode string) (*domain.UseCase, error) {
	items, err := q.rdb.ZRangeWithScores(ctx, q.key(), 0, scanWindow-1).Result()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil // queue is empty – not considered an error
	}

	bestIdx, bestScore := -1, 1e9
	var bestUC domain.UseCase

	for i, it := range items {
		raw, ok := it.Member.(string)
		if !ok {
			continue // skip corrupted element
		}
		var uc domain.UseCase
		if err := json.Unmarshal([]byte(raw), &uc); err != nil {
			continue // also corrupted – ignore
		}

		score := it.Score
		if config.SameScreenGroup(currentNode, uc.Node) {
			score -= screenBoost
		}

		if score < bestScore {
			bestIdx, bestScore, bestUC = i, score, uc
		}
	}

	if bestIdx == -1 {
		// All elements were corrupted – clear the window just in case.
		_ = q.rdb.ZRem(ctx, q.key(), extractMembers(items)...)
		return nil, fmt.Errorf("no decodable use cases in queue window")
	}

	// Remove the selected element by original Member.
	if err := q.rdb.ZRem(ctx, q.key(), items[bestIdx].Member).Err(); err != nil {
		return nil, err
	}

	return &bestUC, nil
}

// extractMembers collects Member values from redis.Z array for utility
// removal of corrupted records.
func extractMembers(zs []redis.Z) []interface{} {
	out := make([]interface{}, 0, len(zs))
	for _, z := range zs {
		out = append(out, z.Member)
	}
	return out
}

// Peek returns the highest priority UseCase without removal (useful for analysis)
func (q *Queue) Peek(ctx context.Context) (*domain.UseCase, error) {
	items, err := q.rdb.ZRangeWithScores(ctx, q.key(), 0, 0).Result()
	if err != nil || len(items) == 0 {
		return nil, err
	}

	var uc domain.UseCase
	if err := json.Unmarshal([]byte(items[0].Member.(string)), &uc); err != nil {
		return nil, err
	}

	return &uc, nil
}

// Len returns the number of tasks in the queue
func (q *Queue) Len(ctx context.Context) (int64, error) {
	return q.rdb.ZCard(ctx, q.key()).Result()
}
