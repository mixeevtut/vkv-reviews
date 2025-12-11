package redis

import (
	"context"
	"fmt"
	"strings"
	"time"
	"vkvideo-reviews/internal/config"

	"github.com/redis/go-redis/v9"
)

type Message struct {
	ID     string
	Values map[string]any
}

type Client struct {
	rdb      *redis.Client
	Stream   string
	Group    string
	Consumer string
	MaxLen   int64
}

func New(c config.Redis) (*Client, error) {
	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Username: c.User,
		Password: c.Pass,
		DB:       c.DB,
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Client{
		rdb:      rdb,
		Stream:   c.Stream,
		Group:    c.Group,
		Consumer: c.Consumer,
		MaxLen:   c.MaxLen,
	}, nil
}

func (c *Client) Close() error { return c.rdb.Close() }

// EnsureGroup — создаёт consumer group, если его ещё нет
func (c *Client) EnsureGroup(ctx context.Context) error {
	err := c.rdb.XGroupCreateMkStream(ctx, c.Stream, c.Group, "$").Err()
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "BUSYGROUP") {
		return nil
	}
	return err
}

// Publish — добавляет сообщение в поток
func (c *Client) Publish(ctx context.Context, fields map[string]any) (string, error) {
	args := &redis.XAddArgs{
		Stream: c.Stream,
		ID:     "*",
		Values: fields,
		Approx: true,
	}
	return c.rdb.XAdd(ctx, args).Result()
}

// Read — читает сообщения для текущего consumer’а (блокирующе)
func (c *Client) Read(ctx context.Context, count int64, blockMs int64) ([]Message, error) {
	res, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    c.Group,
		Consumer: c.Consumer,
		Streams:  []string{c.Stream, ">"},
		Count:    count,
		Block:    time.Duration(blockMs) * time.Millisecond,
	}).Result()

	if err != nil && err != redis.Nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, nil
	}

	var msgs []Message
	for _, s := range res {
		for _, m := range s.Messages {
			msgs = append(msgs, Message{
				ID:     m.ID,
				Values: m.Values,
			})
		}
	}
	return msgs, nil
}

// Ack — подтверждает обработку сообщения
func (c *Client) Ack(ctx context.Context, id string) error {
	return c.rdb.XAck(ctx, c.Stream, c.Group, id).Err()
}

func (c *Client) GetClient() *redis.Client {
	return c.rdb
}
