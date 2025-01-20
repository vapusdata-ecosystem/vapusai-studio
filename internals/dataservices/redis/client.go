package redis

import (
	"context"

	redis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
)

type RedisInterface struct{}

type RedisStore struct {
	Client           *redis.Client
	IsFullRedisStack bool
	logger           zerolog.Logger
}

var (
	defaultJsonPath    = "$"
	defaultJsonPathPtr = "$."
)

// Function to connect with redis server and generate a client
func NewRedisStore(ctx context.Context, conf *Redis, l zerolog.Logger) (*RedisStore, error) {
	cl := redis.NewClient(
		&redis.Options{
			Addr:     conf.URL,
			Username: conf.Username,
			Password: conf.Password,
			DB:       conf.Db,
		})

	// Check if the client is nil
	if cl == nil {
		return nil, dmerrors.DMError(ErrRedisConnection, nil)
	}

	// Ping the redis server to check the connection
	_, err := cl.Ping(ctx).Result()
	if err != nil {
		return nil, dmerrors.DMError(ErrRedisPing, err)
	}

	return &RedisStore{
		Client:           cl,
		IsFullRedisStack: conf.IsFullRedisStack,
		logger:           l,
	}, nil
}

func (r *RedisStore) Close() {
	r.Client.Close()
}

func (r *RedisStore) getPath(path string) string {
	if path == utils.EMPTYSTR {
		return defaultJsonPath
	}
	return defaultJsonPathPtr + path
}

func (r *RedisStore) WrtiteData(ctx context.Context, key, path string, value interface{}) (interface{}, error) {

	if r.IsFullRedisStack {
		return r.writeJsonData(ctx, key, path, value)
	}
	return r.WriteKV(ctx, key, utils.AStructToAString(value))
}

func (r *RedisStore) ReadData(ctx context.Context, key, path string) (interface{}, error) {
	if r.IsFullRedisStack {
		return r.readJsonData(ctx, key, path)
	}
	return r.ReadKV(ctx, key)
}

func (r *RedisStore) DeleteData(ctx context.Context, key, path string) (interface{}, error) {
	if r.IsFullRedisStack {
		return r.deleteJsonData(ctx, key, path)
	}
	return r.DeleteKey(ctx, key)
}

func (r *RedisStore) KeyExists(ctx context.Context, key string) (bool, error) {
	return r.keyExists(ctx, key)
}
