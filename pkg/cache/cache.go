package cache

import (
	"context"
	"errors"
	"time"
)

const (
	TwoWeeks   = time.Hour * 24 * 14
	OneDay     = time.Hour * 24
	OneHour    = time.Hour * 1
	TenMinutes = time.Minute * 10
)

var (
	ErrNil       = errors.New("cache: nil")
	ErrNilResult = errors.New("cache: value is null")
)

type Cache interface {
	Close(ctx context.Context) error
	Get(ctx context.Context, key string, out interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
