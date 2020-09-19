package testutil

import (
	"encoding/json"
	"fmt"
)

type FakeStorage struct {
	ints  map[string]int
	bytes map[string][]byte
}

func NewFakeStorage() FakeStorage {
	return FakeStorage{
		ints:  make(map[string]int),
		bytes: make(map[string][]byte),
	}
}

func (s FakeStorage) GetInt(key string, default_ *int) (int, error) {
	i, ok := s.ints[key]
	if !ok {
		if default_ != nil {
			return *default_, nil
		}
		return 0, fmt.Errorf("missing value at %q and no default", key)
	}
	return i, nil
}

func (s FakeStorage) SetInt(key string, x int) error {
	s.ints[key] = x
	return nil
}

func (s FakeStorage) UpdateInts(f func(vals []int) error, keys ...string) ([]int, error) {
	vals := make([]int, len(keys))
	for i, k := range keys {
		v, ok := s.ints[k]
		if !ok {
			v = 0
		}
		vals[i] = v
	}
	if err := f(vals); err != nil {
		return vals, err
	}
	for i, k := range keys {
		s.ints[k] = vals[i]
	}
	return vals, nil
}

func (s FakeStorage) UpdateJSON(key string, v interface{}, f func() error) error {
	_, _ = s.GetJSON(key, v)

	f()

	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.bytes[key] = bs
	return nil
}

func (s FakeStorage) GetJSON(key string, v interface{}) (bool, error) {
	bs, ok := s.bytes[key]
	if !ok {
		return false, nil
	}

	return true, json.Unmarshal(bs, v)
}
