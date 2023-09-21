package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/orm/crud"

	"github.com/omalloc/kratos-layout/internal/biz"
)

type greeterRepo struct {
	crud.CRUD[biz.Greeter]
	data *Data
	log  *log.Helper
}

// SelectBySomeField implements biz.GreeterRepo.
func (*greeterRepo) SelectBySomeField(ctx context.Context, field string) (*biz.Greeter, error) {
	panic("unimplemented")
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		CRUD: crud.New[biz.Greeter](data.db),
		data: data,
		log:  log.NewHelper(logger),
	}
}
