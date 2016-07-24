package redis

import (
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

// RedisStore is the interface for underlying Redis persistence store.
type RedisStore interface {
	GetConnection() redigo.Conn
}

type redigoStore struct {
	pool *redigo.Pool
}

// NewStore creates a new Redis store.
func NewStore(url string) RedisStore {
	// Build redigo connection.
	pool := &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &redigoStore{
		pool: pool,
	}
}

// GetConnection returns a redigo connection to interact with Redis store.
func (rs *redigoStore) GetConnection() redigo.Conn {
	return rs.pool.Get()
}
