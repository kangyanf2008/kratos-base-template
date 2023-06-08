package data

import (
	"context"

	"kratos-base-template/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type ProviderRepo struct {
	data *Data
	log  *log.Helper
}

func NewProviderRepo(data *Data, logger log.Logger) biz.ProviderBaseInterface {
	return &ProviderRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *ProviderRepo) Save(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *ProviderRepo) Update(ctx context.Context, g *biz.User) (*biz.User, error) {
	return g, nil
}

func (r *ProviderRepo) Get(context.Context, int64) (*biz.User, error) {
	return nil, nil
}

func (r *ProviderRepo) ListUser(context.Context) ([]*biz.User, error) {
	return nil, nil
}
