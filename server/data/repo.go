package data

import (
	"context"

	"github.com/iyabchen/go-react-kv/server/model"
)

// PairRepo specifies interface to fetch data.
type PairRepo interface {
	GetOne(ctx context.Context, id string) (*model.Pair, error)
	GetAll(ctx context.Context) ([]*model.Pair, error)
	DeleteOne(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
	Create(ctx context.Context, p *model.Pair) error
	Update(ctx context.Context, id string, key string, value string) error
}
