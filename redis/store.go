package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Store struct {
	r                  *redis.Client
	MaxConflictRetries int
}

const defaultMaxConflictRetries = 1000

func (s Store) GetInt(key string, default_ *int) (int, error) {
	i, err := s.r.Get(key).Int()
	if err == redis.Nil && default_ != nil {
		return *default_, nil
	}
	return i, err
}

// func (s Store) Inc(key string, delta int) (int, error) {
// }

func (s Store) IncWithin(key string, delta, max int) (int, bool, error) {
	var ErrLimit = fmt.Errorf("limit reached")

	v, updated, err := s.update(key, func(key string, v *redis.StringCmd) (interface{}, error) {
		i, err := v.Int()
		if err == redis.Nil {
			err = nil
		}
		if err != nil {
			return nil, err
		}
		if i+delta <= max {
			return i + delta, nil
		}

		return nil, ErrLimit
	})

	if err == ErrLimit {
		i, ok := v.(int)
		if !ok {
			return 0, updated, fmt.Errorf("unexpected value for %q: not int", key)
		}
		return i, updated, nil
	}

	return 0, updated, err
}

func (s Store) UpdateInts(f func(vals []int) error, keys ...string) error {
	// Transactional function.
	txf := func(tx *redis.Tx) error {
		vals := make([]int, len(keys))
		for i, k := range keys {
			v, err := tx.Get(k).Int()
			if err != nil && err != redis.Nil {
				return err
			}
			vals[i] = v
		}

		if err := f(vals); err != nil {
			return err
		}

		// Operation is commited only if the watched keys remain unchanged.
		_, err := tx.TxPipelined(func(pipe redis.Pipeliner) error {
			for i := range keys {
				if err := pipe.Set(keys[i], vals[i], 0).Err(); err != nil {
					return err
				}
			}
			return nil
		})
		return err
	}

	maxRetries := s.MaxConflictRetries
	if maxRetries == 0 {
		maxRetries = defaultMaxConflictRetries
	}
	for i := 0; i < maxRetries; i++ {
		err := s.r.Watch(txf, keys...)
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}

	return redis.TxFailedErr
}

type UpdateFunc = func(key string, v *redis.StringCmd) (interface{}, error)

// Transactional update.
func (s Store) update(key string, f UpdateFunc) (interface{}, bool, error) {
	var result interface{}

	// Transactional function.
	txf := func(tx *redis.Tx) error {
		n := tx.Get(key)

		v, err := f(key, n)
		if err != nil {
			return err
		}
		result = v

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(func(pipe redis.Pipeliner) error {
			return pipe.Set(key, v, 0).Err()
		})
		return err
	}

	maxRetries := s.MaxConflictRetries
	if maxRetries == 0 {
		maxRetries = defaultMaxConflictRetries
	}
	for i := 0; i < maxRetries; i++ {
		err := s.r.Watch(txf, key)
		if err == nil {
			// Success.
			return result, true, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return nil, false, err
	}

	return nil, false, redis.TxFailedErr
}
