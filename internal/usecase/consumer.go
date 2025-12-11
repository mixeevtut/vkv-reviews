package usecase

import (
	"context"
	"log"
	"strconv"
	"time"
	pkgredis "vkvideo-reviews/pkg/redis"
)

type Consumer struct {
	redis *pkgredis.Client
	uc    RatingRecalcUseCase
}

func NewConsumer(rc *pkgredis.Client, uc RatingRecalcUseCase) *Consumer {
	return &Consumer{
		redis: rc,
		uc:    uc,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	if err := c.redis.EnsureGroup(ctx); err != nil {
		return err
	}
	log.Printf("[consumer] ready: stream=%s group=%s consumer=%s", c.redis.Stream, c.redis.Group, c.redis.Consumer)

	for {
		select {
		case <-ctx.Done():
			log.Println("[consumer] stopping: context canceled")
			return ctx.Err()
		default:
		}

		msgs, err := c.redis.Read(ctx, 64, 3000)
		if err != nil {
			log.Println("[consumer] read:", err)
			time.Sleep(time.Second)
			continue
		}

		if len(msgs) == 0 {
			continue
		}

		for _, m := range msgs {
			msg, ok := m.Values["movie_id"].(string)
			if !ok {
				log.Printf("[consumer] missing movie_id in mesage id=%s", m.ID)
				continue
			}

			movieId, err := strconv.Atoi(msg)
			if err != nil {
				log.Printf("[consumer] bad movie_id=%q in msg id=%s: %v", msg, m.ID, err)
				continue
			}

			log.Printf("[consumer] handle msg id=%s, get movie_id=%d", m.ID, movieId)
			if err := c.uc.RecalcReview(ctx, movieId); err != nil {
				log.Println("[consumer] usecase:", err)
				continue
			}

			if err := c.redis.Ack(ctx, m.ID); err != nil {
				log.Println("[consumer] ack:", err)
			}
		}
	}
}
