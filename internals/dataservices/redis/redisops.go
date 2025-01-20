package redis

import (
	"context"

	dmerrors "github.com/vapusdata-oss/aistudio/core/errors"
	utils "github.com/vapusdata-oss/aistudio/core/utils"
)

func (r *RedisStore) WriteKV(ctx context.Context, key, value string) (interface{}, error) {
	res, err := r.Client.Do(ctx, "SET", key, value).Result()
	if err != nil {
		return utils.EMPTYSTR, dmerrors.DMError(ErrRedisWrite, err)
	}
	return res, nil
}

func (r *RedisStore) keyExists(ctx context.Context, key string) (bool, error) {
	res, err := r.Client.Do(ctx, "Exists", key).Bool()
	if !res {
		return false, dmerrors.DMError(ErrRedisKeyNotExists, err)
	}
	return true, nil
}

func (r *RedisStore) ReadKV(ctx context.Context, key string) (interface{}, error) {
	res, err := r.Client.Do(ctx, "GET", key).Result()
	if err != nil {
		return utils.EMPTYSTR, dmerrors.DMError(ErrRedisRead, err)
	}
	return res, nil
}

func (r *RedisStore) DeleteKey(ctx context.Context, key string) (interface{}, error) {
	res, err := r.Client.Do(ctx, "DEL", key).Result()
	if err != nil {
		return utils.EMPTYSTR, dmerrors.DMError(ErrRedisDelete, err)
	}
	return res, nil
}

func (r *RedisStore) writeJsonData(ctx context.Context, key, path string, value interface{}) (interface{}, error) {
	path = r.getPath(path)
	result := r.Client.JSONSet(ctx, key, path, value)
	if result.Err() != nil {
		return utils.EMPTYSTR, dmerrors.DMError(ErrRedisWrite, result.Err())
	}
	return result.Val(), nil
}

func (r *RedisStore) readJsonData(ctx context.Context, key, path string) (interface{}, error) {
	path = r.getPath(path)
	result := r.Client.JSONGet(ctx, key, path)
	if result.Err() != nil {
		return utils.EMPTYSTR, dmerrors.DMError(ErrRedisRead, result.Err())
	}
	return result.Result()
}

func (r *RedisStore) deleteJsonData(ctx context.Context, key, path string) (interface{}, error) {
	path = r.getPath(path)
	res, err := r.Client.Do(ctx, "JSON.DEL", key, path).Result()
	if err != nil {
		return utils.EMPTYSTR, dmerrors.DMError(ErrRedisDelete, err)
	}
	return res, nil
}

func (r *RedisStore) BFExists(ctx context.Context, key, item string) (interface{}, error) {
	return r.Client.Do(ctx, "BF.EXISTS", key, item).Result()
}

func (r *RedisStore) BFAdd(ctx context.Context, key, item string) (interface{}, error) {
	return r.Client.Do(ctx, "BF.EXISTS", key, item).Result()
}
