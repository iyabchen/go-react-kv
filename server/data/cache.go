package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/iyabchen/go-react-kv/server/model"
)

// Mem stores pair in cache.
type Mem struct {
	mtx   sync.Mutex
	cache map[string]*model.Pair
}

// NewMem creates a Mem instance.
func NewMem() (*Mem, error) {
	return &Mem{
		cache: make(map[string]*model.Pair),
	}, nil
}

// GetOne implements interface.
func (m *Mem) GetOne(ctx context.Context, id string) (*model.Pair, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.cache[id]; !ok {
		return nil, fmt.Errorf("id %s not exist", id)
	}
	return m.cache[id], nil
}

// GetAll implements interface.
func (m *Mem) GetAll(context.Context) ([]*model.Pair, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	var arr []*model.Pair
	for _, p := range m.cache {
		arr = append(arr, p)
	}
	return arr, nil
}

// DeleteAll implements interface.
func (m *Mem) DeleteAll(context.Context) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.cache = make(map[string]*model.Pair)
	return nil
}

// DeleteOne implements interface.
func (m *Mem) DeleteOne(ctx context.Context, id string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.cache[id]; !ok {
		return fmt.Errorf("id %s not exist", id)
	}
	delete(m.cache, id)
	return nil
}

// Create implements interface.
func (m *Mem) Create(ctx context.Context, p *model.Pair) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.cache[p.ID] = p
	return nil
}

// Update implements interface.
func (m *Mem) Update(ctx context.Context, id string, key string, value string) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	if _, ok := m.cache[id]; !ok {
		return fmt.Errorf("id %s not exist", id)
	}
	m.cache[id].Key = key
	m.cache[id].Value = value
	return nil
}
