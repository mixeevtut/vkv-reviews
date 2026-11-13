package usecase

import (
	"context"
	"log"
	"strconv"
	"time"
	"vkvideo-reviews/internal/repository"
	pkgredis "vkvideo-reviews/pkg/redis"
)

type Producer struct {
	repo  repository.RecalcEventRepository
	redis *pkgredis.Client
}

func NewProducer(repo repository.RecalcEventRepository, rc *pkgredis.Client) *Producer {
	return &Producer{
		repo:  repo,
		redis: rc,
	}
}

func (p *Producer) Run(ctx context.Context, pollEvery time.Duration, batch int) {
	t := time.NewTicker(pollEvery)
	defer t.Stop()

	log.Printf("[producer] ready: stream=%s pollEvery=%s batch=%d", p.redis.Stream, pollEvery, batch)

	for {
		select {
		case <-ctx.Done():
			log.Println("[producer] stopping: context canceled")
			return
		case <-t.C:
			evs, err := p.repo.FetchUnsent(ctx, batch)
			if err != nil {
				log.Println("[producer] fetch:", err)
				continue
			}

			for _, e := range evs {
				id, err := p.redis.Publish(ctx, map[string]any{
					"movie_id": strconv.FormatInt(e.MovieID, 10),
				})
				if err != nil {
					log.Println("[producer] publish:", err)
					continue
				}
				if err := p.repo.MarkSent(ctx, e.ID); err != nil {
					log.Println("[producer] mark sent:", err)
					continue
				}
				log.Printf("[producer] sent %s movie_id=%d\n", id, e.MovieID)
			}
		}
	}
}
