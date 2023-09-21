package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/orm"
	"github.com/omalloc/contrib/kratos/orm/crud"
	"github.com/omalloc/contrib/protobuf"

	v1 "github.com/omalloc/kratos-layout/api/helloworld/v1"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	ID    int64 `gorm:"column:id;primaryKey;autoIncrement"`
	Hello string

	orm.DBModel
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	crud.CRUD[Greeter]

	SelectBySomeField(ctx context.Context, field string) (*Greeter, error)
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) error {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.repo.Create(ctx, g)
}

func (uc *GreeterUsecase) List(ctx context.Context, pagination *protobuf.Pagination, _ string) ([]*Greeter, error) {
	return uc.repo.SelectList(ctx, pagination)
}
